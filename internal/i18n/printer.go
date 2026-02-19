package i18n

import (
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// P is the global printer for the application.
var P *message.Printer

//go:generate gotext -srclang=en update -out=catalog.go -lang=en,zh github.com/activebook/auto-wechat github.com/activebook/auto-wechat/internal

func init() {
	// Default to system language or English
	lang := os.Getenv("LANG")
	SetLanguage(lang)

	// test for Chinese
	// SetLanguage("zh")
}

// SetLanguage sets the global printer language.
func SetLanguage(lang string) {
	if strings.Contains(lang, "zh") {
		P = message.NewPrinter(language.Chinese)
	} else {
		P = message.NewPrinter(language.English)
	}
}
