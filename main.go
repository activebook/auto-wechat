package main

import (
	"fmt"
	"os"

	"github.com/activebook/auto-wechat/internal"
	// "github.com/activebook/auto-wechat/test"
)

func main() {
	// test.Test()

	if len(os.Args) < 2 {
		fmt.Println("Usage: auto-wechat <subcommand>")
		fmt.Println("Subcommands:")
		fmt.Println("  scan  - Scan for UI elements and save to settings.yml")
		fmt.Println("  proc  - Run automation process based on settings.yml")
		fmt.Println("  clear - Clear messages that input but not sent")
		return
	}

	subcmd := os.Args[1]
	switch subcmd {
	case "scan":
		internal.RunScanner()
	case "proc":
		internal.RunAutomation()
	case "clear":
		internal.RunClear()
	default:
		fmt.Println("Unknown subcommand:", subcmd)
	}

}
