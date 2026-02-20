package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Point struct {
	X       int `yaml:"x"`
	Y       int `yaml:"y"`
	MarginX int `yaml:"margin_x"`
	MarginY int `yaml:"margin_y"`
	W       int `yaml:"w"`
	H       int `yaml:"h"`
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
	FastMode      bool   `yaml:"fast_mode"` // Whether fast mode or not
	Debug         bool   `yaml:"debug"`     // Whether debug or not
}

func GetAppDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func ReadSettings() (*Settings, error) {
	appDir := GetAppDir()
	settingsPath := filepath.Join(appDir, "cfgs", "settings.yml")

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s. Please copy cfgs/settings-template.yml to cfgs/settings.yml and configure it", settingsPath)
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %v", err)
	}

	var settings Settings
	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		// Log the error but return the partial settings as it was implemented previously
		return &settings, fmt.Errorf("error unmarshaling settings: %v", err)
	}
	return &settings, nil
}

func SaveSettings(settings *Settings) error {
	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("error marshaling settings: %v", err)
	}

	appDir := GetAppDir()
	cfgsDir := filepath.Join(appDir, "cfgs")
	settingsPath := filepath.Join(cfgsDir, "settings.yml")

	if _, err := os.Stat(cfgsDir); os.IsNotExist(err) {
		os.MkdirAll(cfgsDir, 0755)
	}

	err = os.WriteFile(settingsPath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing setting file: %v", err)
	}
	return nil
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
