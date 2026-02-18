package internal

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
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

func LeftClick(msg string) {
	robotgo.Click("left")
	if settings.Debug {
		fmt.Printf("Left Click: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func RightClick(msg string) {
	robotgo.Click("right")
	robotgo.MilliSleep(robotgo.MouseSleep) // because right click need time to response
	if settings.Debug {
		fmt.Printf("Right Click: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressEsc(msg string) {
	robotgo.KeyTap("esc")
	if settings.Debug {
		fmt.Printf("Press Esc: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressEnter(msg string) {
	robotgo.KeyTap("enter")
	if settings.Debug {
		fmt.Printf("Press Enter: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressUp(msg string) {
	robotgo.KeyTap("up")
	if settings.Debug {
		fmt.Printf("Press Up: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressDown(msg string) {
	robotgo.KeyTap("down")
	if settings.Debug {
		fmt.Printf("Press Down: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func GetClipboard() (string, error) {
	return robotgo.ReadAll()
}

func SetClipboard(text string) {
	robotgo.WriteAll(text)
}

// Bug on macos
// func PasteClipboard() {
// 	err := robotgo.CmdV()
// 	if err != nil {
// 		fmt.Println("Error pasting clipboard:", err)
// 	}
// }

func PasteClipboard(msg string) {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`)
		cmd.Run()
	} else {
		robotgo.CmdV()
	}
	robotgo.MilliSleep(settings.Interval) // because paste need time to response
	if settings.Debug {
		fmt.Printf("Paste Clipboard: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func MoveTo(x, y int, msg string) {
	robotgo.Move(x, y)
	if settings.Debug {
		fmt.Printf("Move to (%d, %d): %s\n", x, y, msg)
		robotgo.MilliSleep(1000)
	}

}

func MoveRelative(x, y int, msg string) {
	robotgo.MoveRelative(x, y)
	if settings.Debug {
		fmt.Printf("Move relative (%d, %d): %s\n", x, y, msg)
		robotgo.MilliSleep(1000)
	}
}

func ActivateApp(msg string) error {
	err := robotgo.ActiveName(settings.AppTitle)
	if settings.Debug {
		fmt.Printf("Active App: %s\n", msg)
	}
	return err
}

func HookExit() {
	fmt.Println("--- Please press ctrl + shift + q to stop ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("Stopped.")
		hook.End()
		os.Exit(0)
	})

	s := hook.Start()
	hook.Process(s)
}

func GetAppWindowBounds() (x, y, width, height int, err error) {
	pids, err := robotgo.FindIds(settings.AppTitle)
	if err != nil || len(pids) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("%s not found", settings.AppTitle)
	}

	x, y, width, height = robotgo.GetBounds(pids[0])
	return x, y, width, height, nil
}

const (
	MsgWeChat        = "WeChat"
	MsgContacts      = "Contacts Button"
	MsgMore          = "More Button"
	MsgLastContact   = "Last Contact"
	MsgMessagesBar   = "Messages Bar"
	MsgMessagesPopup = "Messages Popup"
	MsgMessagesInput = "Messages Input Area"
	MsgOutOfScreen   = "Out of Screen"
)

// --- Automation Logic ---

func RunAutomation() {
	// Read Settings
	settings = ReadSettings()
	if settings.Debug {
		fmt.Printf("Loaded settings: %+v\n", settings)
	}

	// Check if imgs folder exists
	appDir := GetAppDir()
	imgsDir := filepath.Join(appDir, "imgs")
	if _, err := os.Stat(imgsDir); os.IsNotExist(err) {
		fmt.Printf("Images directory not found at: %s. Please make sure the 'imgs' folder is present.\n", imgsDir)
		return
	}

	// Activate App
	err := ActivateApp(MsgWeChat)
	if err != nil {
		fmt.Println("active app error: ", err)
		fmt.Println("Please open WeChat app first.")
		return
	}

	// Hook exit message
	HookExit()

	fmt.Println("Starting automation sequence...")

	messages_sent = 0
	for messages_sent < settings.MaxCount {
		GotoLastContact()
		end := GotoNextContact()
		if end {
			fmt.Println("Reach the end of contact list.")
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
	MoveTo(settings.Contacts.X+settings.Contacts.MarginX, settings.Contacts.Y+settings.Contacts.MarginY, MsgContacts)
	LeftClick(MsgContacts) // Left Click

	// Go to More
	MoveTo(settings.More.X+settings.More.MarginX, settings.More.Y+settings.More.MarginY, MsgMore)
	LeftClick(MsgLastContact)  // Left Click
	RightClick(MsgLastContact) // Right Click

	// Press Esc
	PressEsc(MsgLastContact)
	LeftClick(MsgLastContact) // Left Click

	// After doing this, we can now get the focus on the contact list

	// Here we should clip out a small image based on current point
	// For further checking whether the same contact is delt already or not
	if settings.CheckEnd {
		captureContact()
	}
}

func GotoNextContact() bool {
	end := false
	if messages_sent <= 0 {
		return end
	}

	// Press Down
	PressDown(MsgLastContact)

	// Check if contact changed
	if settings.CheckEnd {
		if !checkContactChanged() {
			fmt.Println("Contact did not change (End of list)")
			end = true
		}
	}
	return end
}

func SendMessage() {
	// Go to More
	MoveTo(settings.More.X+settings.More.MarginX, settings.More.Y+settings.More.MarginY, MsgMore)
	RightClick(MsgLastContact) // Right Click

	// Go to Messages
	MoveRelative(settings.MessagesPopup.MarginX, settings.MessagesPopup.MarginY, MsgMessagesPopup)
	LeftClick(MsgMessagesPopup)

	// Go to MessagesBar
	MoveTo(settings.MessagesBar.X+settings.MessagesBar.MarginX, settings.MessagesBar.Y+settings.MessagesBar.MarginY, MsgMessagesBar)
	LeftClick(MsgMessagesInput)

	// Paste clipboard
	PasteClipboard(MsgMessagesInput)

	// Press Enter
	if settings.AutoSend {
		PressEnter(MsgMessagesInput)
	}
}

// Helpers for checks

func normalizeImage(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

func captureAppImage() image.Image {
	ax, ay, aw, ah, err := GetAppWindowBounds()
	if err != nil {
		fmt.Println("get app window bounds error: ", err)
		return nil
	}
	MoveTo(-100, -100, MsgOutOfScreen) // hide cursor
	captured, err := robotgo.CaptureImg(ax, ay, aw, ah)
	if err != nil || captured == nil {
		fmt.Printf("captureAppImage: capture failed: %v\n", err)
		return nil
	}
	// Must save and reload to get correct image
	// robotgo.Save(captured, "imgs/app.png")
	// img, _, _ := robotgo.DecodeImg("imgs/app.png")

	// Bugfix: the raw image returned in memory has a different color format
	// (typically BGRA or some platform-specific byte order)
	// than what gcv.FindImg expects (usually RGBA or standard Go image.RGBA).
	img := normalizeImage(captured)
	return img
}

func captureContact() {
	x, y := robotgo.Location()
	// Capture a small area around the current mouse position (Contact list item)
	// Assuming 100x40 is enough to see the contact name/avatar change
	// Adjust offset to center the capture or ensuring it hits the list item
	MoveTo(-100, -100, MsgOutOfScreen) // hide cursor
	captured, err := robotgo.CaptureImg(x-50, y-20, 100, 40)
	if err != nil || captured == nil {
		fmt.Printf("captureContact: capture failed: %v\n", err)
		return
	}
	// Must save and reload to get correct image
	// robotgo.Save(captured, "imgs/last_contact.png")
	// lastContactImg, _, _ = robotgo.DecodeImg("imgs/last_contact.png")

	// Bugfix: the raw image returned in memory has a different color format
	// (typically BGRA or some platform-specific byte order)
	// than what gcv.FindImg expects (usually RGBA or standard Go image.RGBA).
	lastContactImg = normalizeImage(captured)

	// Get App Image and find contact inside it
	appImg := captureAppImage()
	_, best_match, _, best_points := gcv.FindImg(lastContactImg, appImg)
	if best_match > 0.9 {
		lastContactPos = best_points
		if settings.Debug {
			fmt.Printf("Last contact pos: %v [match:%0.3f]\n", lastContactPos, best_match)
		}
	} else {
		fmt.Printf("Last contact not found: [match:%0.3f]\n", best_match)
	}
}

func checkContactChanged() bool {
	if lastContactImg == nil {
		return true
	}

	// Using gcv.FindImg with high threshold to check equality
	// If FindImg returns a very high match, they are the same.
	appImg := captureAppImage()
	_, bestMatch, _, best_points := gcv.FindImg(lastContactImg, appImg)

	// Since we are capturing the exact same region, identical images should match near 1.0 (or > 0.99)
	// Use 0.95 to account for minor rendering artifacts, or 0.99 for strictness.
	isSame := bestMatch > 0.9

	if isSame && best_points == lastContactPos {
		return false // No change
	}
	return true // Changed
}
