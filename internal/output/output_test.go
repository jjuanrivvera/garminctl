package output

import (
	"bytes"
	"strings"
	"testing"
)

type sample struct {
	Weight float64        `json:"weight"`
	Name   string         `json:"name"`
	Nested map[string]any `json:"nested"`
}

func TestRenderAllFormats(t *testing.T) {
	v := sample{Weight: 72.5, Name: "juan", Nested: map[string]any{"a": 1}}
	for _, format := range []string{"json", "yaml", "csv", "table", "JSON"} {
		var b bytes.Buffer
		if err := Render(&b, format, v); err != nil {
			t.Fatalf("%s: %v", format, err)
		}
		if !strings.Contains(b.String(), "juan") {
			t.Errorf("%s output missing field name: %q", format, b.String())
		}
	}
}

func TestRenderNonObjectWrapped(t *testing.T) {
	var b bytes.Buffer
	if err := Render(&b, "table", []int{1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "value") {
		t.Errorf("array should render under value: %q", b.String())
	}
}

func TestSanitizeTerminal(t *testing.T) {
	if got := SanitizeTerminal("plain text"); got != "plain text" {
		t.Errorf("plain text changed: %q", got)
	}
	dirty := "before\x1b[31mRED\x1b[0mafter"
	got := SanitizeTerminal(dirty)
	if strings.ContainsRune(got, 0x1b) {
		t.Errorf("CSI escape not stripped: %q", got)
	}
	for _, want := range []string{"before", "RED", "after"} {
		if !strings.Contains(got, want) {
			t.Errorf("content lost (%q): %q", want, got)
		}
	}
	osc := "x\x1b]0;evil title\x07y"
	if strings.ContainsRune(SanitizeTerminal(osc), 0x1b) {
		t.Errorf("OSC title sequence not stripped: %q", SanitizeTerminal(osc))
	}
}

func TestRenderCollapsesMultiline(t *testing.T) {
	var b bytes.Buffer
	if err := Render(&b, "table", map[string]any{"note": "line1\nline2"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "line1 line2") {
		t.Errorf("multiline value not collapsed: %q", b.String())
	}
}
