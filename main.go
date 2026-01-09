// main.go
package main

import (
	"ecrypto/cmd"
	"ecrypto/ui"
	"os"
)

func main() {
	// If no CLI args, run interactive menu
	if len(os.Args) == 1 {
		if err := ui.RunInteractiveMenu(); err != nil {
			os.Exit(1)
		}
		return
	}

	// Otherwise, use traditional CLI mode
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}