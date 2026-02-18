package internal

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/go-vgo/robotgo"
)

var (
	settings      Settings
	messages_sent int
)

// --- Utility Functions ---

func LeftClick() {
	robotgo.Click("left")
}

func RightClick() {
	robotgo.Click("right")
}

func PressEsc() {
	robotgo.KeyTap("esc")
}

func PressEnter() {
	robotgo.KeyTap("enter")
}

func PressUp() {
	robotgo.KeyTap("up")
}

func PressDown() {
	robotgo.KeyTap("down")
}

func GetClipboard() (string, error) {
	return robotgo.ReadAll()
}

func SetClipboard(text string) {
	robotgo.WriteAll(text)
}

// func PasteClipboard() {
// 	err := robotgo.CmdV()
// 	if err != nil {
// 		fmt.Println("Error pasting clipboard:", err)
// 	}
// }

func PasteClipboard() {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`)
		cmd.Run()
	} else {
		robotgo.CmdV()
	}
}

func MoveTo(x, y int) {
	robotgo.Move(x, y)
}

func MoveRelative(x, y int) {
	robotgo.MoveRelative(x, y)
}

// --- Automation Logic ---

func RunAutomation() {
	robotgo.MouseSleep = 300
	// Read Settings
	settings := ReadSettings()
	fmt.Printf("Loaded settings: %+v\n", settings)

	fmt.Println("Starting automation sequence...")

	messages_sent = 0
	for messages_sent < settings.MaxCount {
		GotoLastContact()
		GotoNextContact()
		SendMessage()
		messages_sent += 1
	}

	fmt.Println("Automation sequence completed.")
}

func GotoLastContact() {

	// Go to Contacts
	fmt.Println("Moving to Contacts...")
	MoveTo(settings.Contacts.X+settings.Contacts.MarginX, settings.Contacts.Y+settings.Contacts.MarginY)
	LeftClick() // Left Click

	// Go to More
	fmt.Println("Moving to More...")
	MoveTo(settings.More.X+settings.More.MarginX, settings.More.Y+settings.More.MarginY)
	LeftClick()  // Left Click
	RightClick() // Right Click

	// Press Esc
	fmt.Println("Press Esc...")
	PressEsc()  // Press Esc
	LeftClick() // Left Click

}

func GotoNextContact() {
	if messages_sent <= 0 {
		return
	}

	// Press Down
	fmt.Println("Press Down...")
	PressDown()
}

func SendMessage() {
	// Go to More
	fmt.Println("Moving to More...")
	MoveTo(settings.More.X+settings.More.MarginX, settings.More.Y+settings.More.MarginY)
	RightClick() // Right Click

	// Go to Messages
	fmt.Println("Moving to Messages Popup...")
	MoveRelative(settings.MessagesPopup.MarginX, settings.MessagesPopup.MarginY)
	LeftClick()

	// Go to MessagesBar
	fmt.Println("Moving to MessagesBar...")
	MoveTo(settings.MessagesBar.X+settings.MessagesBar.MarginX, settings.MessagesBar.Y+settings.MessagesBar.MarginY)
	LeftClick()

	// Paste clipboard
	fmt.Println("Paste clipboard...")
	// text, err := GetClipboard()
	// if err != nil {
	// 	fmt.Println("Error getting clipboard:", err)
	// 	return
	// }
	// fmt.Println("Clipboard content:\n", text)
	// robotgo.Type(text)
	PasteClipboard()

	// Press Enter
	if settings.AutoSend {
		fmt.Println("Press Enter...")
		PressEnter()
	}
}
