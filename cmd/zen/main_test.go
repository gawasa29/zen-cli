package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"zen-cli/internal/zencli"
)

func TestOptionsFromArgsDefault(t *testing.T) {
	parsed, err := optionsFromArgs(nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandRun {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
	if parsed.dryRun {
		t.Fatal("expected dryRun to be false")
	}
	if parsed.configPathSet {
		t.Fatal("expected configPathSet to be false")
	}
}

func TestOptionsFromArgsListSubcommand(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"list"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandList {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
}

func TestOptionsFromArgsHelpShortFlag(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"-h"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandHelp {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
}

func TestOptionsFromArgsHelpCommand(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"help"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandHelp {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
}

func TestOptionsFromArgsHelpCommandWithTopic(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"help", "add"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandHelp {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
	if parsed.helpTopic != "add" {
		t.Fatalf("unexpected helpTopic: %q", parsed.helpTopic)
	}
}

func TestOptionsFromArgsAddSubcommand(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"add", "Ghostty", "Visual Studio Code,Arc"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandAdd {
		t.Fatalf("unexpected command: %s", parsed.command)
	}

	want := []string{"Ghostty", "Visual Studio Code", "Arc"}
	if !reflect.DeepEqual(parsed.commandApps, want) {
		t.Fatalf("unexpected apps: got %v want %v", parsed.commandApps, want)
	}
}

func TestOptionsFromArgsAddSubcommandWithoutQuotes(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"add", "Visual", "Studio", "Code"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandAdd {
		t.Fatalf("unexpected command: %s", parsed.command)
	}

	want := []string{"Visual Studio Code"}
	if !reflect.DeepEqual(parsed.commandApps, want) {
		t.Fatalf("unexpected apps: got %v want %v", parsed.commandApps, want)
	}
}

func TestOptionsFromArgsAddSubcommandHelp(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"add", "-h"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandHelp {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
	if parsed.helpTopic != "add" {
		t.Fatalf("unexpected helpTopic: %q", parsed.helpTopic)
	}
}

func TestOptionsFromArgsAddRequiresApp(t *testing.T) {
	_, err := optionsFromArgs([]string{"add"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOptionsFromArgsRemoveSubcommand(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"remove", "Ghostty,Arc"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandRemove {
		t.Fatalf("unexpected command: %s", parsed.command)
	}

	want := []string{"Ghostty", "Arc"}
	if !reflect.DeepEqual(parsed.commandApps, want) {
		t.Fatalf("unexpected apps: got %v want %v", parsed.commandApps, want)
	}
}

func TestOptionsFromArgsRemoveSubcommandWithoutQuotes(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"remove", "Visual", "Studio", "Code"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandRemove {
		t.Fatalf("unexpected command: %s", parsed.command)
	}

	want := []string{"Visual Studio Code"}
	if !reflect.DeepEqual(parsed.commandApps, want) {
		t.Fatalf("unexpected apps: got %v want %v", parsed.commandApps, want)
	}
}

func TestOptionsFromArgsRemoveSubcommandHelp(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"remove", "--help"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandHelp {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
	if parsed.helpTopic != "remove" {
		t.Fatalf("unexpected helpTopic: %q", parsed.helpTopic)
	}
}

func TestOptionsFromArgsRemoveRequiresApp(t *testing.T) {
	_, err := optionsFromArgs([]string{"remove"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOptionsFromArgsAllowFlag(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"--allow", "Ghostty, Visual Studio Code"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := []string{"Ghostty", "Visual Studio Code"}
	if !reflect.DeepEqual(parsed.options.AllowedApps, want) {
		t.Fatalf("unexpected AllowedApps: got %v want %v", parsed.options.AllowedApps, want)
	}
}

func TestOptionsFromArgsAllowOnlyFlag(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"--allow-only", "--allow", "Ghostty"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !parsed.options.ReplaceDefaultAllowed {
		t.Fatal("expected ReplaceDefaultAllowed to be true")
	}
	if !parsed.allowOnlySet {
		t.Fatal("expected allowOnlySet to be true")
	}
}

func TestOptionsFromArgsDisallowFlag(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"--disallow", "Ghostty, Finder"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := []string{"Ghostty", "Finder"}
	if !reflect.DeepEqual(parsed.options.DisallowedApps, want) {
		t.Fatalf("unexpected DisallowedApps: got %v want %v", parsed.options.DisallowedApps, want)
	}
}

func TestOptionsFromArgsListFlag(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"--list"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandList {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
}

func TestOptionsFromArgsDryRunAndConfig(t *testing.T) {
	parsed, err := optionsFromArgs([]string{"--dry-run", "--config", "/tmp/custom.json"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if parsed.command != commandRun {
		t.Fatalf("unexpected command: %s", parsed.command)
	}
	if !parsed.dryRun {
		t.Fatal("expected dryRun to be true")
	}
	if !parsed.configPathSet {
		t.Fatal("expected configPathSet to be true")
	}
	if parsed.configPath != "/tmp/custom.json" {
		t.Fatalf("unexpected configPath: %q", parsed.configPath)
	}
}

func TestOptionsFromArgsListAndDryRunConflict(t *testing.T) {
	_, err := optionsFromArgs([]string{"--list", "--dry-run"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultConfigPathUsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-config")

	got, err := defaultConfigPath()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := filepath.Join("/tmp/xdg-config", "zen-cli", "config.json")
	if got != want {
		t.Fatalf("unexpected path: got %q want %q", got, want)
	}
}

func TestDefaultConfigPathFallsBackToHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	got, err := defaultConfigPath()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := filepath.Join(home, ".config", "zen-cli", "config.json")
	if got != want {
		t.Fatalf("unexpected path: got %q want %q", got, want)
	}
}

func TestValidateOptionsAllowOnlyRequiresAllow(t *testing.T) {
	err := validateOptions(zencli.Options{ReplaceDefaultAllowed: true})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMergeOptionsWithConfigBase(t *testing.T) {
	base := zencli.Options{
		AllowedApps:           []string{"Arc"},
		DisallowedApps:        []string{"Slack"},
		ReplaceDefaultAllowed: true,
	}
	cli := zencli.Options{
		AllowedApps:           []string{"Ghostty"},
		DisallowedApps:        []string{"Arc"},
		ReplaceDefaultAllowed: false,
	}

	got := mergeOptions(base, cli, false)
	want := zencli.Options{
		AllowedApps:           []string{"Arc", "Ghostty"},
		DisallowedApps:        []string{"Slack", "Arc"},
		ReplaceDefaultAllowed: true,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected merged options: got %+v want %+v", got, want)
	}
}

func TestMergeOptionsAllowOnlyOverride(t *testing.T) {
	base := zencli.Options{ReplaceDefaultAllowed: true}
	cli := zencli.Options{ReplaceDefaultAllowed: false}

	got := mergeOptions(base, cli, true)
	if got.ReplaceDefaultAllowed {
		t.Fatal("expected ReplaceDefaultAllowed to be false")
	}
}

func TestLoadOptionsFromConfigOptionalMissing(t *testing.T) {
	got, err := loadOptionsFromConfig(filepath.Join(t.TempDir(), "missing.json"), false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !reflect.DeepEqual(got, zencli.Options{}) {
		t.Fatalf("unexpected options: %+v", got)
	}
}

func TestLoadOptionsFromConfigRequiredMissing(t *testing.T) {
	_, err := loadOptionsFromConfig(filepath.Join(t.TempDir(), "missing.json"), true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadOptionsFromConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	body := `{
  "replaceDefaultAllowed": true,
  "allowedApps": ["Ghostty", "Visual Studio Code"],
  "disallowedApps": ["Slack"]
}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	got, err := loadOptionsFromConfig(path, true)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := zencli.Options{
		ReplaceDefaultAllowed: true,
		AllowedApps:           []string{"Ghostty", "Visual Studio Code"},
		DisallowedApps:        []string{"Slack"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected options: got %+v want %+v", got, want)
	}
}

func TestSaveOptionsToConfigRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	opts := zencli.Options{
		ReplaceDefaultAllowed: true,
		AllowedApps:           []string{"Ghostty"},
		DisallowedApps:        []string{"Slack"},
	}

	if err := saveOptionsToConfig(path, opts); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	got, err := loadOptionsFromConfig(path, true)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !reflect.DeepEqual(got, opts) {
		t.Fatalf("unexpected options: got %+v want %+v", got, opts)
	}
}

func TestAddAllowedApps(t *testing.T) {
	base := zencli.Options{
		AllowedApps:    []string{"Ghostty"},
		DisallowedApps: []string{"Arc", "Slack"},
	}

	got := addAllowedApps(base, []string{"Arc", "Visual Studio Code"})
	want := zencli.Options{
		AllowedApps:    []string{"Ghostty", "Arc", "Visual Studio Code"},
		DisallowedApps: []string{"Slack"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected options: got %+v want %+v", got, want)
	}
}

func TestRemoveAllowedApps(t *testing.T) {
	base := zencli.Options{
		AllowedApps:    []string{"Ghostty", "Visual Studio Code"},
		DisallowedApps: []string{"Slack"},
	}

	got := removeAllowedApps(base, []string{"Ghostty", "Arc"})
	want := zencli.Options{
		AllowedApps:    []string{"Visual Studio Code"},
		DisallowedApps: []string{"Slack", "Ghostty", "Arc"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected options: got %+v want %+v", got, want)
	}
}

func TestPrintAllowedApps(t *testing.T) {
	var out bytes.Buffer
	printAllowedApps(&out, []string{"Terminal", "Ghostty"})

	got := out.String()
	want := "zen-cli allowed apps:\n- Terminal\n- Ghostty\n"
	if got != want {
		t.Fatalf("unexpected output: got %q want %q", got, want)
	}
}

func TestPrintDryRunTargets(t *testing.T) {
	var out bytes.Buffer
	printDryRunTargets(&out, []string{"Safari", "Slack"})

	got := out.String()
	want := "zen-cli dry-run targets:\n- Safari\n- Slack\n"
	if got != want {
		t.Fatalf("unexpected output: got %q want %q", got, want)
	}
}

func TestPrintDryRunTargetsEmpty(t *testing.T) {
	var out bytes.Buffer
	printDryRunTargets(&out, nil)

	got := out.String()
	want := "zen-cli dry-run: no target apps would be closed.\n"
	if got != want {
		t.Fatalf("unexpected output: got %q want %q", got, want)
	}
}

func TestPrintHelpRoot(t *testing.T) {
	var out bytes.Buffer
	printHelp(&out, "")

	got := out.String()
	if !strings.Contains(got, "Usage:") {
		t.Fatalf("unexpected output: %q", got)
	}
	if !strings.Contains(got, "zen add APP_NAME") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestPrintHelpAdd(t *testing.T) {
	var out bytes.Buffer
	printHelp(&out, "add")

	got := out.String()
	if !strings.Contains(got, "Usage: zen add APP_NAME") {
		t.Fatalf("unexpected output: %q", got)
	}
}
