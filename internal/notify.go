package internal

import (
	_ "embed"
	"time"

	"github.com/gen2brain/beeep"
)

const ()

func Notify(title, content string) {
	// err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	// if err != nil {
	// 	i18n.P.Printf("Notify error: %v\n", err)
	// }
	beeep.Alert(title, content, "")
	// Optional: add a short delay to ensure the sound plays before the program exits (especially on Windows)
	time.Sleep(time.Millisecond * 200)
}

func Test_notify() {
	beeep.AppName = "MyApp"

	// Send notification with embedded icon
	if err := beeep.Notify("Hello", "This is a desktop notification!", ""); err != nil {
		panic(err)
	}

	// Play a beep
	if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
		panic(err)
	}

	// Alert = notification + beep
	if err := beeep.Alert("Done!", "Task completed.", "icon.png"); err != nil {
		panic(err)
	}
}
