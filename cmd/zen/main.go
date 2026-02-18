package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"zenswitch/internal/zenswitch"
)

type zenCommand string

const (
	commandRun    zenCommand = "run"
	commandList   zenCommand = "list"
	commandAdd    zenCommand = "add"
	commandRemove zenCommand = "remove"
	commandHelp   zenCommand = "help"
)

type parsedArgs struct {
	command       zenCommand
	commandApps   []string
	helpTopic     string
	options       zenswitch.Options
	dryRun        bool
	configPath    string
	configPathSet bool
	allowOnlySet  bool
}

type configFile struct {
	AllowedApps           []string `json:"allowedApps"`
	DisallowedApps        []string `json:"disallowedApps"`
	ReplaceDefaultAllowed bool     `json:"replaceDefaultAllowed"`
}

func main() {
	parsed, err := optionsFromArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
		os.Exit(1)
	}

	if parsed.command == commandHelp {
		printHelp(os.Stdout, parsed.helpTopic)
		return
	}

	configPath := parsed.configPath
	if !parsed.configPathSet {
		configPath, err = defaultConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
			os.Exit(1)
		}
	}

	if parsed.command == commandAdd || parsed.command == commandRemove {
		configOpts, err := loadOptionsFromConfig(configPath, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
			os.Exit(1)
		}

		var updated zenswitch.Options
		if parsed.command == commandAdd {
			updated = addAllowedApps(configOpts, parsed.commandApps)
		} else {
			updated = removeAllowedApps(configOpts, parsed.commandApps)
		}

		if err := saveOptionsToConfig(configPath, updated); err != nil {
			fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "ZenSwitch config updated.")
		printAllowedApps(os.Stdout, zenswitch.EffectiveAllowedApps(updated))
		return
	}

	configOpts, err := loadOptionsFromConfig(configPath, parsed.configPathSet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
		os.Exit(1)
	}

	opts := mergeOptions(configOpts, parsed.options, parsed.allowOnlySet)
	if err := validateOptions(opts); err != nil {
		fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
		os.Exit(1)
	}

	if parsed.command == commandList {
		printAllowedApps(os.Stdout, zenswitch.EffectiveAllowedApps(opts))
		return
	}

	if parsed.dryRun {
		targets, err := zenswitch.PreviewWithOptions(zenswitch.OSExecutor{}, opts)
		if err != nil {
			if errors.Is(err, zenswitch.ErrUnsupportedOS) {
				fmt.Fprintln(os.Stderr, "ZenSwitch is macOS-only.")
				os.Exit(2)
			}
			fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
			os.Exit(1)
		}
		printDryRunTargets(os.Stdout, targets)
		return
	}

	killed, err := zenswitch.ExecuteWithOptions(zenswitch.OSExecutor{}, opts)
	if err != nil {
		if errors.Is(err, zenswitch.ErrUnsupportedOS) {
			fmt.Fprintln(os.Stderr, "ZenSwitch is macOS-only.")
			os.Exit(2)
		}
		fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
		os.Exit(1)
	}

	if len(killed) == 0 {
		fmt.Println("ZenSwitch: no target apps were running.")
		return
	}

	fmt.Println("ZenSwitch closed apps:")
	for _, app := range killed {
		fmt.Printf("- %s\n", app)
	}
}

func optionsFromArgs(args []string) (parsedArgs, error) {
	if len(args) > 0 {
		first := args[0]
		if isHelpToken(first) {
			return parsedArgs{command: commandHelp}, nil
		}
		if first == string(commandHelp) {
			if len(args) > 2 {
				return parsedArgs{}, errors.New("zen help accepts at most one command name")
			}
			topic := ""
			if len(args) == 2 {
				topic = strings.ToLower(strings.TrimSpace(args[1]))
			}
			return parsedArgs{command: commandHelp, helpTopic: topic}, nil
		}

		switch first {
		case string(commandList):
			if len(args) == 2 && isHelpToken(args[1]) {
				return parsedArgs{command: commandHelp, helpTopic: string(commandList)}, nil
			}
			if len(args) != 1 {
				return parsedArgs{}, errors.New("zen list does not accept extra arguments")
			}
			return parsedArgs{command: commandList}, nil
		case string(commandAdd):
			if len(args) == 2 && isHelpToken(args[1]) {
				return parsedArgs{command: commandHelp, helpTopic: string(commandAdd)}, nil
			}
			apps := parseAppArgs(args[1:])
			if len(apps) == 0 {
				return parsedArgs{}, errors.New("zen add requires at least one app name")
			}
			return parsedArgs{command: commandAdd, commandApps: apps}, nil
		case string(commandRemove):
			if len(args) == 2 && isHelpToken(args[1]) {
				return parsedArgs{command: commandHelp, helpTopic: string(commandRemove)}, nil
			}
			apps := parseAppArgs(args[1:])
			if len(apps) == 0 {
				return parsedArgs{}, errors.New("zen remove requires at least one app name")
			}
			return parsedArgs{command: commandRemove, commandApps: apps}, nil
		}
	}

	fs := flag.NewFlagSet("zen", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	allow := fs.String("allow", "", "comma-separated app names to allow")
	allowOnly := fs.Bool("allow-only", false, "use only apps provided by --allow")
	disallow := fs.String("disallow", "", "comma-separated app names to remove from allow-list")
	list := fs.Bool("list", false, "print effective allow-list and exit")
	dryRun := fs.Bool("dry-run", false, "show target apps and exit without closing")
	config := fs.String("config", "", "path to config JSON file")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return parsedArgs{command: commandHelp}, nil
		}
		return parsedArgs{}, err
	}

	seen := make(map[string]struct{})
	fs.Visit(func(f *flag.Flag) {
		seen[f.Name] = struct{}{}
	})

	command := commandRun
	if *list {
		command = commandList
	}
	if *list && *dryRun {
		return parsedArgs{}, errors.New("--list and --dry-run cannot be used together")
	}

	return parsedArgs{
		command: command,
		options: zenswitch.Options{
			AllowedApps:           parseAllowApps(*allow),
			DisallowedApps:        parseAllowApps(*disallow),
			ReplaceDefaultAllowed: *allowOnly,
		},
		dryRun:        *dryRun,
		configPath:    strings.TrimSpace(*config),
		configPathSet: hasFlag(seen, "config"),
		allowOnlySet:  hasFlag(seen, "allow-only"),
	}, nil
}

func hasFlag(flags map[string]struct{}, name string) bool {
	_, ok := flags[name]
	return ok
}

func isHelpToken(arg string) bool {
	return arg == "-h" || arg == "--help"
}

func defaultConfigPath() (string, error) {
	if xdg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); xdg != "" {
		return filepath.Join(xdg, "zenswitch", "config.json"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	return filepath.Join(home, ".config", "zenswitch", "config.json"), nil
}

func loadOptionsFromConfig(path string, required bool) (zenswitch.Options, error) {
	if strings.TrimSpace(path) == "" {
		return zenswitch.Options{}, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && !required {
			return zenswitch.Options{}, nil
		}
		return zenswitch.Options{}, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg configFile
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return zenswitch.Options{}, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return zenswitch.Options{
		AllowedApps:           cfg.AllowedApps,
		DisallowedApps:        cfg.DisallowedApps,
		ReplaceDefaultAllowed: cfg.ReplaceDefaultAllowed,
	}, nil
}

func saveOptionsToConfig(path string, opts zenswitch.Options) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("config path is empty")
	}

	cfg := configFile{
		AllowedApps:           append([]string{}, opts.AllowedApps...),
		DisallowedApps:        append([]string{}, opts.DisallowedApps...),
		ReplaceDefaultAllowed: opts.ReplaceDefaultAllowed,
	}

	body, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	body = append(body, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}
	return nil
}

func mergeOptions(base zenswitch.Options, cli zenswitch.Options, allowOnlySet bool) zenswitch.Options {
	merged := zenswitch.Options{
		AllowedApps:           append([]string{}, base.AllowedApps...),
		DisallowedApps:        append([]string{}, base.DisallowedApps...),
		ReplaceDefaultAllowed: base.ReplaceDefaultAllowed,
	}

	merged.AllowedApps = append(merged.AllowedApps, cli.AllowedApps...)
	merged.DisallowedApps = append(merged.DisallowedApps, cli.DisallowedApps...)

	if allowOnlySet {
		merged.ReplaceDefaultAllowed = cli.ReplaceDefaultAllowed
	}

	return merged
}

func validateOptions(opts zenswitch.Options) error {
	if opts.ReplaceDefaultAllowed && len(opts.AllowedApps) == 0 {
		return errors.New("--allow-only requires allow apps in CLI or config")
	}
	return nil
}

func addAllowedApps(opts zenswitch.Options, apps []string) zenswitch.Options {
	opts.AllowedApps = mergeAppLists(opts.AllowedApps, apps)
	opts.DisallowedApps = removeFromAppList(opts.DisallowedApps, apps)
	return opts
}

func removeAllowedApps(opts zenswitch.Options, apps []string) zenswitch.Options {
	opts.AllowedApps = removeFromAppList(opts.AllowedApps, apps)
	opts.DisallowedApps = mergeAppLists(opts.DisallowedApps, apps)
	return opts
}

func mergeAppLists(base []string, incoming []string) []string {
	result := make([]string, 0, len(base)+len(incoming))
	seen := make(map[string]struct{}, len(base)+len(incoming))

	appendUnique := func(name string) {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			return
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}

	for _, app := range base {
		appendUnique(app)
	}
	for _, app := range incoming {
		appendUnique(app)
	}

	return result
}

func removeFromAppList(base []string, toRemove []string) []string {
	if len(toRemove) == 0 {
		return mergeAppLists(nil, base)
	}

	removed := make(map[string]struct{}, len(toRemove))
	for _, app := range toRemove {
		trimmed := strings.TrimSpace(app)
		if trimmed == "" {
			continue
		}
		removed[strings.ToLower(trimmed)] = struct{}{}
	}

	filtered := make([]string, 0, len(base))
	for _, app := range mergeAppLists(nil, base) {
		if _, blocked := removed[strings.ToLower(strings.TrimSpace(app))]; blocked {
			continue
		}
		filtered = append(filtered, app)
	}
	return filtered
}

func printAllowedApps(out io.Writer, apps []string) {
	if len(apps) == 0 {
		fmt.Fprintln(out, "ZenSwitch allowed apps: (none)")
		return
	}

	fmt.Fprintln(out, "ZenSwitch allowed apps:")
	for _, app := range apps {
		fmt.Fprintf(out, "- %s\n", app)
	}
}

func printDryRunTargets(out io.Writer, apps []string) {
	if len(apps) == 0 {
		fmt.Fprintln(out, "ZenSwitch dry-run: no target apps would be closed.")
		return
	}

	fmt.Fprintln(out, "ZenSwitch dry-run targets:")
	for _, app := range apps {
		fmt.Fprintf(out, "- %s\n", app)
	}
}

func printHelp(out io.Writer, topic string) {
	switch topic {
	case string(commandList):
		fmt.Fprintln(out, "Usage: zen list")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Show effective allowed apps without closing applications.")
		return
	case string(commandAdd):
		fmt.Fprintln(out, "Usage: zen add APP_NAME [APP_NAME ...]")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Add app names to the allow-list and persist to config.")
		return
	case string(commandRemove):
		fmt.Fprintln(out, "Usage: zen remove APP_NAME [APP_NAME ...]")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Remove app names from the allow-list and persist to config.")
		return
	}

	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, "  zen")
	fmt.Fprintln(out, "  zen list")
	fmt.Fprintln(out, "  zen add APP_NAME [APP_NAME ...]")
	fmt.Fprintln(out, "  zen remove APP_NAME [APP_NAME ...]")
	fmt.Fprintln(out, "  zen help [list|add|remove]")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Options:")
	fmt.Fprintln(out, "  --dry-run                    Show target apps and exit without closing")
	fmt.Fprintln(out, "  --config PATH               Use a specific config file path")
	fmt.Fprintln(out, "  --allow APP1,APP2           Append allow apps for this run")
	fmt.Fprintln(out, "  --allow-only                Use only explicitly allowed apps")
	fmt.Fprintln(out, "  --disallow APP1,APP2        Remove allow apps for this run")
	fmt.Fprintln(out, "  --list                      List effective allow apps (legacy)")
	fmt.Fprintln(out, "  -h, --help                  Show help")
}

func parseAllowApps(raw string) []string {
	if strings.TrimSpace(raw) == "" {
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

func parseAppArgs(args []string) []string {
	if len(args) == 0 {
		return nil
	}

	for _, arg := range args {
		if strings.Contains(arg, ",") {
			apps := make([]string, 0, len(args))
			for _, part := range args {
				apps = append(apps, parseAllowApps(part)...)
			}
			return apps
		}
	}

	joined := strings.TrimSpace(strings.Join(args, " "))
	if joined == "" {
		return nil
	}
	return []string{joined}
}
