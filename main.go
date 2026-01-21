// main.go
package main

import (
	"ecrypto/cmd"
	"ecrypto/gui"
	"ecrypto/ui"
	"flag"
	"fmt"
	"log"
	"os"
)

var Version = "dev"

func main() {
    // Print version if requested
    if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
        fmt.Printf("ecrypto %s\n", Version)
        return
    }

    // Check for --serve flag (API server mode for GUI)
    serveFlag := flag.Bool("serve", false, "Run HTTP API server for GUI")
    portFlag := flag.Int("port", 8765, "API server port")
    flag.Parse()

    if *serveFlag {
        server := gui.NewServer(*portFlag)
        log.Fatal(server.Start())
        return
    }

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