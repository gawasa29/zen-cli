package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestOptionsFromArgsDefault(t *testing.T) {
	opts, listOnly, err := optionsFromArgs(nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if listOnly {
		t.Fatal("expected listOnly to be false")
	}

	if opts.ReplaceDefaultAllowed {
		t.Fatal("expected ReplaceDefaultAllowed to be false")
	}

	if len(opts.AllowedApps) != 0 {
		t.Fatalf("expected empty AllowedApps, got %v", opts.AllowedApps)
	}
}

func TestOptionsFromArgsAllow(t *testing.T) {
	opts, listOnly, err := optionsFromArgs([]string{"--allow", "Ghostty, Visual Studio Code"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if listOnly {
		t.Fatal("expected listOnly to be false")
	}

	want := []string{"Ghostty", "Visual Studio Code"}
	if !reflect.DeepEqual(opts.AllowedApps, want) {
		t.Fatalf("unexpected AllowedApps: got %v want %v", opts.AllowedApps, want)
	}
}

func TestOptionsFromArgsAllowOnlyRequiresAllow(t *testing.T) {
	_, _, err := optionsFromArgs([]string{"--allow-only"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOptionsFromArgsAllowOnly(t *testing.T) {
	opts, listOnly, err := optionsFromArgs([]string{"--allow-only", "--allow", "Ghostty"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if listOnly {
		t.Fatal("expected listOnly to be false")
	}

	if !opts.ReplaceDefaultAllowed {
		t.Fatal("expected ReplaceDefaultAllowed to be true")
	}
}

func TestOptionsFromArgsDisallow(t *testing.T) {
	opts, listOnly, err := optionsFromArgs([]string{"--disallow", "Ghostty, Finder"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if listOnly {
		t.Fatal("expected listOnly to be false")
	}

	want := []string{"Ghostty", "Finder"}
	if !reflect.DeepEqual(opts.DisallowedApps, want) {
		t.Fatalf("unexpected DisallowedApps: got %v want %v", opts.DisallowedApps, want)
	}
}

func TestOptionsFromArgsList(t *testing.T) {
	_, listOnly, err := optionsFromArgs([]string{"--list"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !listOnly {
		t.Fatal("expected listOnly to be true")
	}
}

func TestPrintAllowedApps(t *testing.T) {
	var out bytes.Buffer
	printAllowedApps(&out, []string{"Terminal", "Ghostty"})

	got := out.String()
	want := "ZenSwitch allowed apps:\n- Terminal\n- Ghostty\n"
	if got != want {
		t.Fatalf("unexpected output: got %q want %q", got, want)
	}
}
