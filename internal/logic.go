package internal

import (
	"fmt"
	"image"
	"os/exec"
	"runtime"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

var (
	settings       Settings
	messages_sent  int
	lastContactImg image.Image
	lastContactPos image.Point
)

// --- Utility Functions ---

func init() {
	robotgo.MouseSleep = 100 // add this at the start of your program
	robotgo.KeySleep = 50
}

func LeftClick() {
	robotgo.Click("left")
	// robotgo.MilliSleep(50)
}

func RightClick() {
	robotgo.Click("right")
	// robotgo.MilliSleep(50)
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
	robotgo.MilliSleep(100)
}

func MoveTo(x, y int) {
	robotgo.Move(x, y)
	// robotgo.MilliSleep(100)
}

func MoveRelative(x, y int) {
	robotgo.MoveRelative(x, y)
	// robotgo.MilliSleep(100)
}

func ActiveApp() error {
	return robotgo.ActiveName(settings.AppTitle)

}

// --- Automation Logic ---

func RunAutomation() {
	// Read Settings
	settings = ReadSettings()
	fmt.Printf("Loaded settings: %+v\n", settings)

	err := ActiveApp()
	if err != nil {
		fmt.Println("active app error: ", err)
		return
	}

	fmt.Println("Starting automation sequence...")

	messages_sent = 0
	for messages_sent < settings.MaxCount {
		GotoLastContact()
		end := GotoNextContact()
		if end {
			fmt.Println("End of contact list or checks failed.")
			break
		}

		SendMessage()
		messages_sent += 1
		robotgo.MilliSleep(settings.Interval)
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

	// After doing this, we can now get the focus on the contact list

	// Here we should clip out a small image based on current point
	// For further checking whether the same contact is delt already or not
	if settings.CheckLastOne {
		captureContact()
	}
}

func GotoNextContact() bool {
	end := false
	if messages_sent <= 0 {
		return end
	}

	// Press Down
	fmt.Println("Press Down...")
	PressDown()

	// Check if state changed
	if settings.CheckLastOne {
		if !checkContactChanged() {
			fmt.Println("State did not change (End of list)")
			end = true
		}
	}
	return end
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
	PasteClipboard()

	// Press Enter
	if settings.AutoSend {
		fmt.Println("Press Enter...")
		PressEnter()
	}
}

// Helpers for checks

func captureContact() {
	x, y := robotgo.Location()
	// Capture a small area around the current mouse position (Contact list item)
	// Assuming 100x40 is enough to see the contact name/avatar change
	// Adjust offset to center the capture or ensuring it hits the list item
	MoveTo(-100, -100) // hide cursor
	lastContactImg, _ := robotgo.CaptureImg(x-50, y-20, 100, 40)
	robotgo.Save(lastContactImg, "imgs/last_one.png")
	img, _ := robotgo.CaptureImg()
	robotgo.Save(img, "imgs/screen.png")
	_, best_match, _, best_points := gcv.FindImg(lastContactImg, img)
	if best_match > 0.8 {
		lastContactPos = best_points
		fmt.Printf("Last contact pos: %v\n", lastContactPos)
	} else {
		fmt.Printf("Last contact not found: %f\n", best_match)
	}
}

func checkContactChanged() bool {
	if lastContactImg == nil {
		return true
	}

	// Using gcv.FindImg with high threshold to check equality
	// If FindImg returns a very high match, they are the same.
	img, _ := robotgo.CaptureImg()
	_, bestMatch, _, best_points := gcv.FindImg(lastContactImg, img)

	// Since we are capturing the exact same region, identical images should match near 1.0 (or > 0.99)
	// Use 0.95 to account for minor rendering artifacts, or 0.99 for strictness.
	isSame := bestMatch > 0.90

	if isSame && best_points == lastContactPos {
		return false // No change
	}
	return true // Changed
}
