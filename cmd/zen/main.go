package main

import (
	"errors"
	"fmt"
	"os"

	"zenswitch/internal/zenswitch"
)

func main() {
	killed, err := zenswitch.Execute(zenswitch.OSExecutor{})
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
