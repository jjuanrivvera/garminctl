package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestHookScript_BashExecution runs the generated guard hook through real bash to verify the
// adversarial cases: obfuscation, path-invoked binaries, the method-gated api hatch, and the
// benign lookalikes that must stay allowed. Gated on a POSIX shell being available.
func TestHookScript_BashExecution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}

	hookFile := filepath.Join(t.TempDir(), "garminctl-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookScript()), 0o755); err != nil { // #nosec G306 -- hook must be executable
		t.Fatalf("write hook: %v", err)
	}

	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  "Bash",
			"tool_input": map[string]any{"command": command},
		})
		return string(b)
	}
	mcpPayload := func(toolName string) string {
		b, _ := json.Marshal(map[string]any{"tool_name": toolName, "tool_input": map[string]any{}})
		return string(b)
	}

	runHook := func(t *testing.T, payload string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Run(); err != nil {
			t.Logf("hook output: %s", out.String())
			t.Fatalf("hook script exited non-zero: %v", err)
		}
		return out.String()
	}
	isDenied := func(output string) bool { return strings.Contains(output, `"permissionDecision":"deny"`) }

	cases := []struct {
		name       string
		payload    string
		wantDenied bool
	}{
		// --- blocked commands ---
		{"auth_logout_denied", bashPayload("garminctl auth logout"), true},
		{"auth_logout_profile_denied", bashPayload("garminctl --profile juan auth logout"), true},
		{"alias_set_denied", bashPayload(`garminctl alias set kill "auth logout"`), true},
		// --- obfuscation ---
		{"quote_split_denied", bashPayload(`garminctl auth log""out`), true},
		{"single_quote_split_denied", bashPayload(`garminctl auth log''out`), true},
		{"backslash_denied", bashPayload(`garminctl auth log\out`), true},
		// --- command position after separators / env prefix ---
		{"after_semicolon_denied", bashPayload("true; garminctl auth logout"), true},
		{"after_pipe_denied", bashPayload("echo hi | garminctl auth logout"), true},
		{"env_prefix_denied", bashPayload("GARMINCTL_PROFILE=juan garminctl auth logout"), true},
		// --- path-invoked binaries ---
		{"relative_path_binary_denied", bashPayload("./bin/garminctl auth logout"), true},
		{"absolute_path_binary_denied", bashPayload("/usr/local/bin/garminctl auth logout"), true},
		// --- raw api hatch gated by HTTP method ---
		{"api_delete_denied", bashPayload("garminctl api /activity-service/activity/123 -X DELETE"), true},
		{"api_post_denied", bashPayload("garminctl api /weight-service/weight --method POST --data {}"), true},
		{"api_put_lowercase_denied", bashPayload("garminctl api /x --method put"), true},
		{"api_profile_before_denied", bashPayload("garminctl --profile juan api /x -X PUT"), true},
		{"api_compound_read_then_delete_denied", bashPayload("garminctl sleep; garminctl api /x -X DELETE"), true},
		// --- raw api reads stay allowed ---
		{"api_get_allowed", bashPayload("garminctl api /usersummary-service/usersummary/daily"), false},
		{"api_explicit_get_allowed", bashPayload("garminctl api /x -X GET"), false},
		{"api_delete_in_value_allowed", bashPayload(`garminctl api /x --data '{"note":"delete later"}'`), false},
		// --- benign lookalikes that must stay allowed ---
		{"sleep_allowed", bashPayload("garminctl sleep"), false},
		{"body_comp_allowed", bashPayload("garminctl body-composition --date 2026-07-10"), false},
		{"auth_status_allowed", bashPayload("garminctl auth status"), false},
		{"cat_file_allowed", bashPayload("cat auth_logout.go"), false},
		{"other_binary_allowed", bashPayload("mygarminctl auth logout"), false},
		{"config_use_with_logout_word_allowed", bashPayload(`garminctl config use "notes about auth logout"`), false},
		// --- MCP branch: garminctl's MCP surface is read-only, so nothing is blocked ---
		{"mcp_stress_allowed", mcpPayload("mcp__garminctl__garmin_stress"), false},
		{"mcp_sleep_allowed", mcpPayload("mcp__garminctl__garmin_sleep"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if denied := isDenied(runHook(t, tc.payload)); denied != tc.wantDenied {
				t.Errorf("want denied=%v, got denied=%v", tc.wantDenied, denied)
			}
		})
	}
}

// TestHookScript_BashExecutionNoJq exercises the no-jq fallback with a STRICT PATH: a bin dir
// holding only the POSIX tools the hook needs, so jq is genuinely unreachable (merely
// prepending an empty dir leaves jq resolvable — the flaw that masks fail-open bugs).
func TestHookScript_BashExecutionNoJq(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("bash hook tests require a POSIX shell; skipping on windows")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH; skipping hook execution tests")
	}

	tmpDir := t.TempDir()
	hookFile := filepath.Join(tmpDir, "garminctl-guard.sh")
	if err := os.WriteFile(hookFile, []byte(hookScript()), 0o755); err != nil { // #nosec G306 -- hook must be executable
		t.Fatalf("write hook: %v", err)
	}

	binDir := filepath.Join(tmpDir, "nojq-bin")
	if err := os.Mkdir(binDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	for _, tool := range []string{"cat", "tr", "grep", "sed", "printf", "env"} {
		p, lerr := exec.LookPath(tool)
		if lerr != nil {
			continue
		}
		if serr := os.Symlink(p, filepath.Join(binDir, tool)); serr != nil {
			t.Fatalf("symlink %s: %v", tool, serr)
		}
	}

	bashPayload := func(command string) string {
		b, _ := json.Marshal(map[string]any{
			"tool_name":  "Bash",
			"tool_input": map[string]any{"command": command},
		})
		return string(b)
	}
	runHookNoJq := func(t *testing.T, payload string) string {
		t.Helper()
		cmd := exec.Command(bash, hookFile)
		cmd.Stdin = strings.NewReader(payload)
		env := make([]string, 0, len(os.Environ()))
		for _, e := range os.Environ() {
			if !strings.HasPrefix(e, "PATH=") {
				env = append(env, e)
			}
		}
		cmd.Env = append(env, "PATH="+binDir)
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Run(); err != nil {
			t.Logf("hook output: %s", out.String())
			t.Fatalf("hook script exited non-zero: %v", err)
		}
		return out.String()
	}
	isDenied := func(output string) bool { return strings.Contains(output, `"permissionDecision":"deny"`) }

	cases := []struct {
		name       string
		payload    string
		wantDenied bool
	}{
		{"nojq_auth_logout_denied", bashPayload("garminctl auth logout"), true},
		{"nojq_obfuscated_logout_denied", bashPayload(`garminctl auth log""out`), true},
		{"nojq_path_binary_denied", bashPayload("/usr/local/bin/garminctl auth logout"), true},
		{"nojq_api_delete_denied", bashPayload("garminctl api /x -X DELETE"), true},
		{"nojq_sleep_allowed", bashPayload("garminctl sleep"), false},
		{"nojq_api_get_allowed", bashPayload("garminctl api /usersummary-service/usersummary/daily"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if denied := isDenied(runHookNoJq(t, tc.payload)); denied != tc.wantDenied {
				t.Errorf("want denied=%v, got denied=%v", tc.wantDenied, denied)
			}
		})
	}
}
