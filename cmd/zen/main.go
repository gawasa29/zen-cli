package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"zenswitch/internal/zenswitch"
)

func main() {
	opts, listOnly, err := optionsFromArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ZenSwitch failed: %v\n", err)
		os.Exit(1)
	}

	if listOnly {
		printAllowedApps(os.Stdout, zenswitch.EffectiveAllowedApps(opts))
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

func optionsFromArgs(args []string) (zenswitch.Options, bool, error) {
	fs := flag.NewFlagSet("zen", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	allow := fs.String("allow", "", "comma-separated app names to allow")
	allowOnly := fs.Bool("allow-only", false, "use only apps provided by --allow")
	disallow := fs.String("disallow", "", "comma-separated app names to remove from allow-list")
	list := fs.Bool("list", false, "print effective allow-list and exit")

	if err := fs.Parse(args); err != nil {
		return zenswitch.Options{}, false, err
	}

	opts := zenswitch.Options{
		AllowedApps:           parseAllowApps(*allow),
		DisallowedApps:        parseAllowApps(*disallow),
		ReplaceDefaultAllowed: *allowOnly,
	}

	if opts.ReplaceDefaultAllowed && len(opts.AllowedApps) == 0 {
		return zenswitch.Options{}, false, errors.New("--allow-only requires --allow")
	}

	return opts, *list, nil
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
