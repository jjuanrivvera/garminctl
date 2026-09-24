package commands

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/zalando/go-keyring"
)

// The payloads below are the shapes Garmin actually returns; the empty ones are from issue #16.
func TestHasData(t *testing.T) {
	cases := []struct {
		name string
		json string
		want bool
	}{{
		name: "unsynced sleep day: nulls and zeros under an echoed date",
		json: `{"dailySleepDTO":{"id":null,"calendarDate":"2026-09-17","sleepStartTimestampGMT":0,
		        "sleepEndTimestampGMT":0,"sleepTimeSeconds":0,"deepSleepSeconds":null}}`,
		want: false,
	}, {
		name: "unsynced heart rate",
		json: `{"userProfilePK":123,"calendarDate":"2026-09-17","restingHeartRate":0,
		        "heartRateValues":null}`,
		want: false,
	}, {
		name: "a slept night",
		json: `{"dailySleepDTO":{"id":42,"calendarDate":"2026-09-17","sleepTimeSeconds":27000}}`,
		want: true,
	}, {
		name: "a resting heart rate of 52",
		json: `{"userProfilePK":123,"calendarDate":"2026-09-17","restingHeartRate":52}`,
		want: true,
	}, {
		name: "empty object",
		json: `{}`,
		want: false,
	}, {
		name: "empty array",
		json: `[]`,
		want: false,
	}, {
		name: "array of readings",
		json: `[{"value":0},{"value":98}]`,
		want: true,
	}, {
		name: "a flag that is set is data",
		json: `{"calendarDate":"2026-09-17","napAlertEnabled":true}`,
		want: true,
	}, {
		name: "nothing but the request echoed back",
		json: `{"calendarDate":"2026-09-17","userProfileId":9,"uuid":"abc"}`,
		want: false,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var v any
			if err := json.Unmarshal([]byte(tc.json), &v); err != nil {
				t.Fatal(err)
			}
			if got := hasData(v); got != tc.want {
				t.Errorf("hasData = %v, want %v", got, tc.want)
			}
		})
	}
}

// A value the JSON round-trip cannot represent must count as data: the point is never to hide
// output, only to label a day that is provably empty.
func TestHasData_UnmarshalableCountsAsData(t *testing.T) {
	if !hasData(make(chan int)) {
		t.Error("a payload that cannot be marshalled must be treated as data")
	}
}

// End to end: the mocked API returns {} for a day, so the read succeeds, the output still
// prints, stderr explains, and the process exits with ExitNoData rather than 0 (issue #16).
func TestDailyRead_EmptyDaySignalsNoData(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	testHTTPClient = mockOK()
	t.Cleanup(func() { testHTTPClient = nil })

	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "me", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}

	stdout, stderr, err := execRoot(t, "--profile", "me", "heart-rate", "--date", "2026-09-17", "-o", "json")
	if err != nil {
		t.Fatalf("an empty day is not an error: %v", err)
	}
	if stdout == "" {
		t.Error("the payload must still be printed — consumers that don't check are unaffected")
	}
	if !strings.Contains(stderr, "no heart-rate data for 2026-09-17") {
		t.Errorf("stderr should explain the empty day, got %q", stderr)
	}
	if !gf.noData {
		t.Error("the empty day must be recorded for the exit status")
	}

	// The exit status is the whole point of the fix, so assert the number a script would see
	// rather than the flag behind it.
	if code := Main(t.Context(), []string{
		"--profile", "me", "heart-rate", "--date", "2026-09-17", "-o", "json",
	}); code != ExitNoData {
		t.Errorf("Main returned %d for an empty day, want ExitNoData (%d)", code, ExitNoData)
	}
}

// And a day with real numbers stays silent and exits 0.
func TestDailyRead_RealDayIsSilent(t *testing.T) {
	keyring.MockInit()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GARMINCTL_PROFILE", "")
	testHTTPClient = mockOK()
	t.Cleanup(func() { testHTTPClient = nil })

	gdir := t.TempDir()
	writeGarthTokens(t, gdir, time.Now().Add(time.Hour).Unix())
	if _, _, err := execRoot(t, "--profile", "me", "auth", "import", "--from", gdir); err != nil {
		t.Fatal(err)
	}

	// `steps` is the one resource the mock answers with real numbers.
	_, stderr, err := execRoot(t, "--profile", "me", "steps", "--date", "2026-07-10", "-o", "json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stderr, "no steps data") {
		t.Errorf("a day with measurements must not be flagged: %q", stderr)
	}
	if gf.noData {
		t.Error("a day with measurements must not set the no-data status")
	}
	if code := Main(t.Context(), []string{
		"--profile", "me", "steps", "--date", "2026-07-10", "-o", "json",
	}); code != 0 {
		t.Errorf("Main returned %d for a day with data, want 0", code)
	}
}
