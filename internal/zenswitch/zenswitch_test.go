package zenswitch

import (
	"errors"
	"reflect"
	"testing"
)

type callResult struct {
	output []byte
	err    error
}

type mockExecutor struct {
	calls   []string
	results map[string]callResult
}

func (m *mockExecutor) Run(name string, args ...string) ([]byte, error) {
	call := name
	for _, arg := range args {
		call += "|" + arg
	}
	m.calls = append(m.calls, call)
	if result, ok := m.results[call]; ok {
		return result.output, result.err
	}
	return nil, nil
}

func TestParseAppList(t *testing.T) {
	input := "Safari, Visual Studio Code, Slack"
	got := parseAppList(input)
	want := []string{"Safari", "Visual Studio Code", "Slack"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected parse result: got %v want %v", got, want)
	}
}

func TestFilterTargets(t *testing.T) {
	running := []string{"Terminal", "Safari", "Slack"}
	allowed := makeAllowedSet([]string{"Terminal"})
	got := filterTargets(running, allowed)
	want := []string{"Safari", "Slack"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected targets: got %v want %v", got, want)
	}
}

func TestResolveAllowedAppsDefaultPlusUser(t *testing.T) {
	opts := Options{
		AllowedApps: []string{"Ghostty", "Visual Studio Code"},
	}

	got := resolveAllowedApps(opts)
	want := []string{
		"Terminal",
		"iTerm2",
		"Ghostty",
		"Finder",
		"Dock",
		"System Settings",
		"Activity Monitor",
		"Visual Studio Code",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected allow-list: got %v want %v", got, want)
	}
}

func TestResolveAllowedAppsUserOnly(t *testing.T) {
	opts := Options{
		AllowedApps:           []string{"Ghostty", "Visual Studio Code"},
		ReplaceDefaultAllowed: true,
	}

	got := resolveAllowedApps(opts)
	want := []string{"Ghostty", "Visual Studio Code"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected allow-list: got %v want %v", got, want)
	}
}

func TestResolveAllowedAppsWithDisallowed(t *testing.T) {
	opts := Options{
		AllowedApps:    []string{"Visual Studio Code"},
		DisallowedApps: []string{"ghostty", "Terminal", "Visual Studio Code"},
	}

	got := resolveAllowedApps(opts)
	want := []string{
		"iTerm2",
		"Finder",
		"Dock",
		"System Settings",
		"Activity Monitor",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected allow-list: got %v want %v", got, want)
	}
}

func TestEffectiveAllowedApps(t *testing.T) {
	opts := Options{
		AllowedApps:    []string{"Visual Studio Code"},
		DisallowedApps: []string{"Ghostty"},
	}

	got := EffectiveAllowedApps(opts)
	want := []string{
		"Terminal",
		"iTerm2",
		"Finder",
		"Dock",
		"System Settings",
		"Activity Monitor",
		"Visual Studio Code",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected allow-list: got %v want %v", got, want)
	}
}

func TestTargetAppsFromRunning(t *testing.T) {
	running := []string{"Safari", "Terminal", "Ghostty", "zen", "Visual Studio Code"}
	opts := Options{
		AllowedApps:    []string{"Visual Studio Code"},
		DisallowedApps: []string{"Ghostty"},
	}

	got := targetAppsFromRunning(running, opts)
	want := []string{"Ghostty", "Safari"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected targets: got %v want %v", got, want)
	}
}

func TestQuitAppPkillNoProcessIsAccepted(t *testing.T) {
	mock := &mockExecutor{
		results: map[string]callResult{
			"pkill|-x|Safari": {err: errors.New("exit status 1")},
		},
	}
	if err := quitApp(mock, "Safari"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestQuitAppPkillFailureWithOutput(t *testing.T) {
	mock := &mockExecutor{
		results: map[string]callResult{
			"pkill|-x|Safari": {output: []byte("permission denied"), err: errors.New("exit status 3")},
		},
	}
	if err := quitApp(mock, "Safari"); err == nil {
		t.Fatal("expected error, got nil")
	}
}
