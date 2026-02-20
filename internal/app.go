package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/activebook/auto-wechat/internal/i18n"
	hook "github.com/robotn/gohook"
)

func init() {
	runtime.LockOSThread()
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

func checkImgsFolder() bool {
	appDir := GetAppDir()
	imgsDir := filepath.Join(appDir, "imgs")
	if _, err := os.Stat(imgsDir); os.IsNotExist(err) {
		i18n.P.Printf("Images directory not found at: %s. Please make sure the 'imgs' folder is present.\n", imgsDir)
		return false
	}
	return true
}

func RunAutomation() {
	settings, err := ReadSettings()
	if err != nil {
		i18n.P.Printf("Failed to load settings: %v\n", err)
		return
	}
	if settings.Debug {
		i18n.P.Printf("Loaded settings: %+v\n", settings)
	}

	if !checkImgsFolder() {
		return
	}

	robot := NewRobot(settings.Debug, settings.Interval)
	err = robot.ActivateApp(settings.AppTitle, MsgWeChat)
	if err != nil {
		i18n.P.Printf("active app error: %v\n", err)
		i18n.P.Printf("Please open WeChat app first.\n")
		return
	}

	HookExit()

	vision := NewVision(robot, settings.AppTitle, settings.Debug)
	wechat := NewWeChat(settings, robot, vision)

	i18n.P.Printf("Starting automation sequence...\n")

	for wechat.GetMessagesSent() < settings.MaxCount {
		wechat.GotoLastContact()
		end := wechat.GotoNextContact()
		if end {
			i18n.P.Printf("Reach the end of contact list.\n")
			break
		}

		wechat.SendMessage()
		wechat.IncrementMessagesSent()
		robot.Wait(MsgMessagesInput)
	}

	i18n.P.Printf("Automation sequence completed.\n")
}

func RunClear() {
	settings, err := ReadSettings()
	if err != nil {
		i18n.P.Printf("Failed to load settings: %v\n", err)
		return
	}
	if settings.Debug {
		i18n.P.Printf("Loaded settings: %+v\n", settings)
	}

	if !checkImgsFolder() {
		return
	}

	robot := NewRobot(settings.Debug, settings.Interval)
	err = robot.ActivateApp(settings.AppTitle, MsgWeChat)
	if err != nil {
		i18n.P.Printf("active app error: %v\n", err)
		i18n.P.Printf("Please open WeChat app first.\n")
		return
	}

	HookExit()

	vision := NewVision(robot, settings.AppTitle, settings.Debug)
	wechat := NewWeChat(settings, robot, vision)

	i18n.P.Printf("Starting automation sequence...\n")

	for wechat.GetMessagesSent() < settings.MaxCount {
		wechat.GotoLastContact()
		end := wechat.GotoNextContact()
		if end {
			i18n.P.Printf("Reach the end of contact list.\n")
			break
		}

		wechat.ClearMessage()
		wechat.IncrementMessagesSent()
	}

	i18n.P.Printf("Automation sequence completed.\n")
}

func RunScanner() {
	settings, err := ReadSettings()
	if err != nil {
		// Fallback to empty settings for scanner if unable to read
		settings = &Settings{
			AppTitle: "WeChat",
		}
	}

	robot := NewRobot(settings.Debug, settings.Interval)
	err = robot.ActivateApp(settings.AppTitle, MsgWeChat)
	if err != nil {
		fmt.Println("active app error: ", err)
		fmt.Println("Please open WeChat app first.")
		return
	}

	vision := NewVision(robot, settings.AppTitle, settings.Debug)
	appDir := GetAppDir()

	// 1. Scan Contacts
	contactsImgs := []string{
		filepath.Join(appDir, "imgs/contacts_light.png"),
		filepath.Join(appDir, "imgs/contacts_light_selected.png"),
		filepath.Join(appDir, "imgs/contacts_dark.png"),
		filepath.Join(appDir, "imgs/contacts_dark_selected.png"),
	}
	if p, err := vision.FindPoint(contactsImgs); err == nil {
		fmt.Printf("Found Contacts at: %v\n", p)
		settings.Contacts.X = p.X
		settings.Contacts.Y = p.Y
	} else {
		fmt.Println("Contacts not found")
	}

	// 2. Scan More
	moreImgs := []string{
		filepath.Join(appDir, "imgs/more_light.png"),
		filepath.Join(appDir, "imgs/more_dark.png"),
	}
	if p, err := vision.FindPoint(moreImgs); err == nil {
		fmt.Printf("Found More at: %v\n", p)
		settings.More.X = p.X
		settings.More.Y = p.Y
	} else {
		fmt.Println("More not found")
	}

	// 3. Scan MessagesBar
	messagesbarImgs := []string{
		filepath.Join(appDir, "imgs/messages_bar_light.png"),
		filepath.Join(appDir, "imgs/messages_bar_dark.png"),
	}
	if p, err := vision.FindPoint(messagesbarImgs); err == nil {
		fmt.Printf("Found MessagesBar at: %v\n", p)
		settings.MessagesBar.X = p.X
		settings.MessagesBar.Y = p.Y
	} else {
		fmt.Println("MessagesBar not found")
	}

	// Save to settings.yml
	SaveSettings(settings)
}
