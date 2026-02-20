package main

import (
	"os"

	"github.com/activebook/auto-wechat/internal"
	"github.com/activebook/auto-wechat/internal/i18n"
	// "github.com/activebook/auto-wechat/test"
)

func main() {
	// test.Test()

	if len(os.Args) < 2 {
		i18n.P.Printf("Usage: auto-wechat <subcommand>\n")
		i18n.P.Printf("Description: auto-wechat is a computer vision-based automation tool for WeChat, designed to bulk-send messages in clipboard to all contacts.\n\n")
		i18n.P.Printf("Subcommands:\n")
		i18n.P.Printf("  scan  - Scan for UI elements and save to settings.yml\n")
		i18n.P.Printf("  proc  - Run automation process based on settings.yml\n")
		i18n.P.Printf("  clear - Clear messages that input but not sent\n")
		return
	}

	subcmd := os.Args[1]
	switch subcmd {
	case "scan":
		internal.RunScan()
	case "proc":
		internal.RunProc()
	case "clear":
		internal.RunClear()
	default:
		i18n.P.Printf("Unknown subcommand: %s\n", subcmd)
	}

}
