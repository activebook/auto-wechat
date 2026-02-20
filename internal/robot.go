package internal

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/activebook/auto-wechat/internal/i18n"
	"github.com/go-vgo/robotgo"
)

// Robot handles OS-level interactions (mouse, keyboard, clipboard)
type Robot struct {
	debug    bool
	interval int
}

func NewRobot(debug bool, interval int) *Robot {
	robotgo.MouseSleep = 100
	robotgo.KeySleep = 50

	return &Robot{
		debug:    debug,
		interval: interval,
	}
}

func (r *Robot) LeftClick(msg string) {
	robotgo.Click("left")
	if r.debug {
		i18n.P.Printf("Left Click: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) RightClick(msg string) {
	robotgo.Click("right")
	robotgo.MilliSleep(robotgo.MouseSleep) // because right click need time to response
	if r.debug {
		i18n.P.Printf("Right Click: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) PressEsc(msg string) {
	robotgo.KeyTap(robotgo.Esc)
	if r.debug {
		i18n.P.Printf("Press Esc: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) PressEnter(msg string) {
	robotgo.KeyTap(robotgo.Enter)
	if r.debug {
		i18n.P.Printf("Press Enter: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) PressCtrlA(msg string) {
	if runtime.GOOS == "darwin" {
		exec.Command("osascript", "-e", `tell application "System Events" to keystroke "a" using {command down}`).Run()
	} else {
		robotgo.KeyTap(robotgo.KeyA, robotgo.CmdCtrl())
	}
	robotgo.MilliSleep(r.interval) // because Ctrl+A need time to response
	if r.debug {
		i18n.P.Printf("Press Ctrl+A: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) PressDelete(msg string) {
	robotgo.KeyTap(robotgo.Delete)
	if r.debug {
		i18n.P.Printf("Press Delete: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) PressUp(msg string) {
	robotgo.KeyTap(robotgo.Up)
	if r.debug {
		i18n.P.Printf("Press Up: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) PressDown(msg string) {
	robotgo.KeyTap(robotgo.Down)
	if r.debug {
		i18n.P.Printf("Press Down: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) GetClipboard() (string, error) {
	return robotgo.ReadAll()
}

func (r *Robot) SetClipboard(text string) {
	robotgo.WriteAll(text)
}

func (r *Robot) PasteClipboard(msg string) {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`)
		cmd.Run()
	} else {
		robotgo.CmdV()
	}
	robotgo.MilliSleep(r.interval) // because paste need time to response
	if r.debug {
		i18n.P.Printf("Paste Clipboard: %s\n", msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) MoveTo(x, y int, msg string) {
	robotgo.Move(x, y)
	if r.debug {
		i18n.P.Printf("Move to (%d, %d): %s\n", x, y, msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) MoveRelative(x, y int, msg string) {
	robotgo.MoveRelative(x, y)
	if r.debug {
		i18n.P.Printf("Move relative (%d, %d): %s\n", x, y, msg)
		robotgo.MilliSleep(1000)
	}
}

func (r *Robot) Wait(msg string) {
	robotgo.MilliSleep(r.interval)
	if r.debug {
		i18n.P.Printf("Wait: %s\n", msg)
	}
}

func (r *Robot) ActivateApp(appTitle string, msg string) error {
	if runtime.GOOS == "darwin" {
		err := robotgo.ActiveName(appTitle)
		if r.debug {
			i18n.P.Printf("Active App: %s\n", msg)
		}
		return err
	}

	pids, err := robotgo.FindIds(appTitle)
	if err != nil || len(pids) == 0 {
		return fmt.Errorf("%s not found", appTitle)
	}

	for _, pid := range pids {
		x, y, width, height := robotgo.GetBounds(pid)
		if width > 0 && height > 0 {
			i18n.P.Printf("Found app %s (pid=%d) at [%d, %d, %d, %d]\n", appTitle, pid, x, y, width, height)
			return robotgo.ActivePid(pid)
		}
	}

	return fmt.Errorf("%s not found", appTitle)
}

func (r *Robot) GetAppWindowBounds(appTitle string) (x, y, width, height int, err error) {
	pids, err := robotgo.FindIds(appTitle)
	if err != nil || len(pids) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("%s not found", appTitle)
	}

	for _, pid := range pids {
		x, y, width, height = robotgo.GetBounds(pid)
		if width > 0 && height > 0 {
			return x, y, width, height, nil
		}
	}

	return 0, 0, 0, 0, fmt.Errorf("%s not found", appTitle)
}
