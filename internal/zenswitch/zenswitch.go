package zenswitch

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

var ErrUnsupportedOS = errors.New("zenswitch supports macOS only")

// defaultAllowedApps is the predefined allow-list required by the product.
var defaultAllowedApps = []string{
	"Terminal",
	"iTerm2",
	"Ghostty",
	"Finder",
	"Dock",
	"System Settings",
	"Activity Monitor",
}

type Executor interface {
	Run(name string, args ...string) ([]byte, error)
}

type OSExecutor struct{}

func (OSExecutor) Run(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

type Options struct {
	AllowedApps           []string
	DisallowedApps        []string
	ReplaceDefaultAllowed bool
}

func Execute(executor Executor) ([]string, error) {
	return ExecuteWithOptions(executor, Options{})
}

func EffectiveAllowedApps(opts Options) []string {
	return resolveAllowedApps(opts)
}

func ExecuteWithOptions(executor Executor, opts Options) ([]string, error) {
	if runtime.GOOS != "darwin" {
		return nil, ErrUnsupportedOS
	}

	running, err := runningAppNames(executor)
	if err != nil {
		return nil, err
	}

	allowed := makeAllowedSet(resolveAllowedApps(opts))
	for _, app := range selfExecutableNames() {
		allowed[strings.ToLower(app)] = struct{}{}
	}

	killed := make([]string, 0, len(running))
	for _, app := range filterTargets(running, allowed) {
		if err := quitApp(executor, app); err != nil {
			return killed, err
		}
		killed = append(killed, app)
	}

	sort.Strings(killed)
	return killed, nil
}

func resolveAllowedApps(opts Options) []string {
	allowed := make([]string, 0, len(defaultAllowedApps)+len(opts.AllowedApps))
	if !opts.ReplaceDefaultAllowed {
		allowed = append(allowed, defaultAllowedApps...)
	}

	seen := make(map[string]struct{}, len(allowed))
	for _, app := range allowed {
		seen[strings.ToLower(strings.TrimSpace(app))] = struct{}{}
	}

	for _, app := range opts.AllowedApps {
		name := strings.TrimSpace(app)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		allowed = append(allowed, name)
	}

	if len(opts.DisallowedApps) == 0 {
		return allowed
	}

	disallowed := make(map[string]struct{}, len(opts.DisallowedApps))
	for _, app := range opts.DisallowedApps {
		name := strings.TrimSpace(app)
		if name == "" {
			continue
		}
		disallowed[strings.ToLower(name)] = struct{}{}
	}

	filtered := make([]string, 0, len(allowed))
	for _, app := range allowed {
		if _, blocked := disallowed[strings.ToLower(strings.TrimSpace(app))]; blocked {
			continue
		}
		filtered = append(filtered, app)
	}

	return filtered
}

func selfExecutableNames() []string {
	return []string{"zen"}
}

func runningAppNames(executor Executor) ([]string, error) {
	script := `tell application "System Events" to get name of every application process whose background only is false`
	out, err := executor.Run("osascript", "-e", script)
	if err != nil {
		return nil, fmt.Errorf("failed to list running apps: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return parseAppList(string(out)), nil
}

func quitApp(executor Executor, appName string) error {
	safeName := strings.ReplaceAll(appName, `"`, `\\\"`)
	quitScript := fmt.Sprintf(`tell application "%s" to quit`, safeName)
	if out, err := executor.Run("osascript", "-e", quitScript); err != nil {
		return fmt.Errorf("failed to quit %s: %w: %s", appName, err, strings.TrimSpace(string(out)))
	}

	if out, err := executor.Run("pkill", "-x", appName); err != nil {
		if strings.TrimSpace(string(out)) == "" {
			// already closed
			return nil
		}
		return fmt.Errorf("failed to force close %s: %w: %s", appName, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func parseAppList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	apps := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		apps = append(apps, name)
	}
	return apps
}

func makeAllowedSet(apps []string) map[string]struct{} {
	result := make(map[string]struct{}, len(apps))
	for _, app := range apps {
		result[strings.ToLower(app)] = struct{}{}
	}
	return result
}

func filterTargets(running []string, allowed map[string]struct{}) []string {
	targets := make([]string, 0, len(running))
	for _, app := range running {
		if _, ok := allowed[strings.ToLower(app)]; ok {
			continue
		}
		targets = append(targets, app)
	}
	return targets
}
