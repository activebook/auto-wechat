package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"gopkg.in/yaml.v3"
)

type Point struct {
	X       int `yaml:"x"`
	Y       int `yaml:"y"`
	MarginX int `yaml:"margin_x"`
	MarginY int `yaml:"margin_y"`
}

type Settings struct {
	Contacts      Point  `yaml:"contacts"`
	More          Point  `yaml:"more"`
	MessagesBar   Point  `yaml:"messages_bar"`
	MessagesPopup Point  `yaml:"messages_popup"`
	AppTitle      string `yaml:"app_title"` // WeChat App Title to activate
	Interval      int    `yaml:"interval"`  // Interval for every contact (ms)
	MaxCount      int    `yaml:"max_count"` // Maximum number of messages to send
	CheckEnd      bool   `yaml:"check_end"` // Whether check the last one is the end of contacts list
	AutoSend      bool   `yaml:"auto_send"` // Whether auto send message or not
	Debug         bool   `yaml:"debug"`     // Whether debug or not
}

func GetAppDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func RunScanner() {
	// Load existing settings if available to preserve margins
	// Note: ReadSettings now terminates if file not found, but for scanner we might want to allow it?
	// Actually scanner creates it, so we should check if we can read it, if not, start empty or default.
	// But ReadSettings handles loading. If it fails, it exits.
	// For Scanner: checking if file exists is good to preserve margins.
	// The requirement Bug 2 says "if cfgs folder and file don't exist, it should prompt error to tell user to get it first".
	// This likely applies to 'proc' mode (RunAutomation).
	// For 'scan', we probably want to allow running even if setttings.yml is partial?
	// But let's stick to the request: prompt error.
	settings = ReadSettings()

	// Activate App
	err := ActivateApp(MsgWeChat)
	if err != nil {
		fmt.Println("active app error: ", err)
		fmt.Println("Please open WeChat app first.")
		return
	}

	appDir := GetAppDir()

	// 1. Scan Contacts
	// Try light and light_selected first, then dark and dark_selected
	contactsImgs := []string{
		filepath.Join(appDir, "imgs/contacts_light.png"),
		filepath.Join(appDir, "imgs/contacts_light_selected.png"),
		filepath.Join(appDir, "imgs/contacts_dark.png"),
		filepath.Join(appDir, "imgs/contacts_dark_selected.png"),
	}
	if p, err := findPoint(contactsImgs); err == nil {
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
	if p, err := findPoint(moreImgs); err == nil {
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
	if p, err := findPoint(messagesbarImgs); err == nil {
		fmt.Printf("Found MessagesBar at: %v\n", p)
		settings.MessagesBar.X = p.X
		settings.MessagesBar.Y = p.Y
	} else {
		fmt.Println("MessagesBar not found")
	}

	// Save to settings.yml
	saveSettings(settings)
}

func findPoint(images []string) (Point, error) {
	// Capture the screen once to search against
	screenshot, _ := robotgo.CaptureImg()

	for _, imgPath := range images {
		// Verify file exists
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			fmt.Printf("Image file does not exist: %s\n", imgPath)
			continue
		}

		// Decode the template image
		template, _, err := robotgo.DecodeImg(imgPath)
		if err != nil {
			fmt.Printf("Error decoding image %s: %v\n", imgPath, err)
			continue
		}

		// Find the image
		// x, y = gcv.FindImg(template, screenshot) returns the top-left coordinate
		_, bestMatch, _, bestPoint := gcv.FindImg(template, screenshot)

		// Check if a good match was found (threshold 0.8 is usually good)
		if bestMatch > 0.8 {
			return Point{X: bestPoint.X, Y: bestPoint.Y}, nil
		}
	}

	return Point{}, fmt.Errorf("image not found in logic")
}

func saveSettings(settings Settings) {
	data, err := yaml.Marshal(&settings)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	appDir := GetAppDir()
	cfgsDir := filepath.Join(appDir, "cfgs")
	settingsPath := filepath.Join(cfgsDir, "settings.yml")

	// Ensure cfgs directory exists
	if _, err := os.Stat(cfgsDir); os.IsNotExist(err) {
		os.Mkdir(cfgsDir, 0755)
	}

	err = os.WriteFile(settingsPath, data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("Settings saved to %s\n", settingsPath)
}

func ReadSettings() Settings {
	var settings Settings
	appDir := GetAppDir()
	settingsPath := filepath.Join(appDir, "cfgs", "settings.yml")

	// Check if file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		log.Fatalf("Config file not found at %s. Please copy cfgs/settings-template.yml to cfgs/settings.yml and configure it.", settingsPath)
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Fatalf("error reading settings: %v", err)
	}

	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		log.Printf("error unmarshaling settings: %v", err)
	}
	return settings
}
