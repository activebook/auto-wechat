package internal

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/activebook/auto-wechat/internal/i18n"

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
	runtime.LockOSThread()

	robotgo.MouseSleep = 100 // add this at the start of your program
	robotgo.KeySleep = 50
}

func LeftClick(msg string) {
	robotgo.Click("left")
	if settings.Debug {
		i18n.P.Printf("Left Click: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func RightClick(msg string) {
	robotgo.Click("right")
	robotgo.MilliSleep(robotgo.MouseSleep) // because right click need time to response
	if settings.Debug {
		i18n.P.Printf("Right Click: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressEsc(msg string) {
	robotgo.KeyTap(robotgo.Esc)
	if settings.Debug {
		i18n.P.Printf("Press Esc: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressEnter(msg string) {
	robotgo.KeyTap(robotgo.Enter)
	if settings.Debug {
		i18n.P.Printf("Press Enter: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressCtrlA(msg string) {
	if runtime.GOOS == "darwin" {
		exec.Command("osascript", "-e",
			`tell application "System Events" to keystroke "a" using {command down}`).Run()
	} else {
		robotgo.KeyTap(robotgo.KeyA, robotgo.CmdCtrl())
	}
	robotgo.MilliSleep(settings.Interval) // because Ctrl+A need time to response
	if settings.Debug {
		i18n.P.Printf("Press Ctrl+A: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressDelete(msg string) {
	robotgo.KeyTap(robotgo.Delete)
	if settings.Debug {
		i18n.P.Printf("Press Delete: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressUp(msg string) {
	robotgo.KeyTap(robotgo.Up)
	if settings.Debug {
		i18n.P.Printf("Press Up: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func PressDown(msg string) {
	robotgo.KeyTap(robotgo.Down)
	if settings.Debug {
		i18n.P.Printf("Press Down: %s\n", msg)
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
// 		i18n.P.Println("Error pasting clipboard:", err)
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
		i18n.P.Printf("Paste Clipboard: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func MoveTo(x, y int, msg string) {
	robotgo.Move(x, y)
	if settings.Debug {
		i18n.P.Printf("Move to (%d, %d): %s\n", x, y, msg)
		robotgo.MilliSleep(1000)
	}

}

func MoveRelative(x, y int, msg string) {
	robotgo.MoveRelative(x, y)
	if settings.Debug {
		i18n.P.Printf("Move relative (%d, %d): %s\n", x, y, msg)
		robotgo.MilliSleep(1000)
	}
}

func ActivateApp(msg string) error {
	if runtime.GOOS == "darwin" {
		err := robotgo.ActiveName(settings.AppTitle)
		if settings.Debug {
			i18n.P.Printf("Active App: %s\n", msg)
		}
		return err
	}

	// Windows platform
	// One title cannot have multi ids
	pids, err := robotgo.FindIds(settings.AppTitle)
	if err != nil || len(pids) == 0 {
		return fmt.Errorf("%s not found", settings.AppTitle)
	}

	for _, pid := range pids {
		// Check whether pid has window
		x, y, width, height := robotgo.GetBounds(pid)
		if width > 0 && height > 0 {
			i18n.P.Printf("Found app %s(%d) at [%d, %d, %d, %d]\n", settings.AppTitle, pid, x, y, width, height)
			return robotgo.ActivePid(pid)
		}
	}

	return fmt.Errorf("%s not found", settings.AppTitle)
}

func HookExit() {
	i18n.P.Printf("--- Please press ctrl + shift + q to stop ---\n")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		i18n.P.Printf("Stopped.\n")
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

	for _, pid := range pids {
		x, y, width, height = robotgo.GetBounds(pid)
		if width > 0 && height > 0 {
			return x, y, width, height, nil
		}
	}

	return 0, 0, 0, 0, fmt.Errorf("%s not found", settings.AppTitle)
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
		i18n.P.Printf("Loaded settings: %+v\n", settings)
	}

	// Check if imgs folder exists
	appDir := GetAppDir()
	imgsDir := filepath.Join(appDir, "imgs")
	if _, err := os.Stat(imgsDir); os.IsNotExist(err) {
		i18n.P.Printf("Images directory not found at: %s. Please make sure the 'imgs' folder is present.\n", imgsDir)
		return
	}

	// Activate App
	err := ActivateApp(MsgWeChat)
	if err != nil {
		i18n.P.Printf("active app error: %v\n", err)
		i18n.P.Printf("Please open WeChat app first.\n")
		return
	}

	// Hook exit message
	HookExit()

	i18n.P.Printf("Starting automation sequence...\n")

	messages_sent = 0
	for messages_sent < settings.MaxCount {
		GotoLastContact()
		end := GotoNextContact()
		if end {
			i18n.P.Printf("Reach the end of contact list.\n")
			break
		}

		SendMessage()
		messages_sent += 1
		robotgo.MilliSleep(settings.Interval)
	}

	i18n.P.Printf("Automation sequence completed.\n")
}

func RunClear() {
	// Read Settings
	settings = ReadSettings()
	if settings.Debug {
		i18n.P.Printf("Loaded settings: %+v\n", settings)
	}

	// Check if imgs folder exists
	appDir := GetAppDir()
	imgsDir := filepath.Join(appDir, "imgs")
	if _, err := os.Stat(imgsDir); os.IsNotExist(err) {
		i18n.P.Printf("Images directory not found at: %s. Please make sure the 'imgs' folder is present.\n", imgsDir)
		return
	}

	// Activate App
	err := ActivateApp(MsgWeChat)
	if err != nil {
		i18n.P.Printf("active app error: %v\n", err)
		i18n.P.Printf("Please open WeChat app first.\n")
		return
	}

	// Hook exit message
	HookExit()

	i18n.P.Printf("Starting automation sequence...\n")

	messages_sent = 0
	for messages_sent < settings.MaxCount {
		GotoLastContact()
		end := GotoNextContact()
		if end {
			i18n.P.Printf("Reach the end of contact list.\n")
			break
		}

		ClearMessage()
		messages_sent += 1
	}

	i18n.P.Printf("Automation sequence completed.\n")
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
			i18n.P.Printf("Contact did not change (End of list)\n")
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

func ClearMessage() {
	// Go to More
	MoveTo(settings.More.X+settings.More.MarginX, settings.More.Y+settings.More.MarginY, MsgMore)
	RightClick(MsgLastContact) // Right Click

	// Go to Messages
	MoveRelative(settings.MessagesPopup.MarginX, settings.MessagesPopup.MarginY, MsgMessagesPopup)
	LeftClick(MsgMessagesPopup)

	// Go to MessagesBar
	MoveTo(settings.MessagesBar.X+settings.MessagesBar.MarginX, settings.MessagesBar.Y+settings.MessagesBar.MarginY, MsgMessagesBar)
	LeftClick(MsgMessagesInput)

	// Press CtrlA
	PressCtrlA(MsgMessagesInput)

	// Press Delete
	PressDelete(MsgMessagesInput)
}

// Helpers for checks

func normalizeImage(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

func captureAppImage() image.Image {
	foundApp := true
	ax, ay, aw, ah, err := GetAppWindowBounds()
	if err != nil {
		i18n.P.Printf("get app window bounds error: %v\n", err)
		foundApp = false
	}
	MoveTo(-100, -100, MsgOutOfScreen) // hide cursor

	var captured image.Image
	// capture app image
	if foundApp {
		captured, err = robotgo.CaptureImg(ax, ay, aw, ah)
		if err != nil {
			i18n.P.Printf("captureAppImage[%d, %d, %d, %d]: capture failed: %v\n", ax, ay, aw, ah, err)
		}
	}

	// cannot find app then capture screen
	if captured == nil {
		captured, err = robotgo.CaptureImg()
		if err != nil || captured == nil {
			i18n.P.Printf("capture screen error: %v\n", err)
			return nil
		}
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
		i18n.P.Printf("captureContact: capture failed: %v\n", err)
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
	if appImg == nil {
		return
	}
	_, best_match, _, best_points := gcv.FindImg(lastContactImg, appImg)
	if best_match > 0.9 {
		lastContactPos = best_points
		if settings.Debug {
			i18n.P.Printf("Last contact pos: %v [match:%0.3f]\n", lastContactPos, best_match)
		}
	} else {
		i18n.P.Printf("Last contact not found: [match:%0.3f]\n", best_match)
	}
}

// Check whether the last contact is changed or not
func checkContactChanged() bool {
	if lastContactImg == nil {
		return true
	}

	// Using gcv.FindImg with high threshold to check equality
	// If FindImg returns a very high match, they are the same.
	appImg := captureAppImage()
	if appImg == nil {
		return true
	}
	_, bestMatch, _, best_points := gcv.FindImg(lastContactImg, appImg)

	// Since we are capturing the exact same region, identical images should match near 1.0 (or > 0.99)
	// Use 0.95 to account for minor rendering artifacts, or 0.99 for strictness.
	isSame := bestMatch > 0.9

	if isSame && best_points == lastContactPos {
		return false // No change
	}
	return true // Changed
}
