package internal

import (
	"image"

	"github.com/activebook/auto-wechat/internal/i18n"
	"github.com/go-vgo/robotgo"
)

// WeChat encapsulates automation logic tailored specifically for the WeChat application
type WeChat struct {
	config *Settings
	robot  *Robot
	vision *Vision

	messagesSent   int
	lastContactImg image.Image
	lastContactPos image.Point
}

func NewWeChat(config *Settings, robot *Robot, vision *Vision) *WeChat {
	return &WeChat{
		config: config,
		robot:  robot,
		vision: vision,
	}
}

func (w *WeChat) GotoLastContact() {
	// Go to Contacts
	w.robot.MoveTo(w.config.Contacts.X+w.config.Contacts.MarginX, w.config.Contacts.Y+w.config.Contacts.MarginY, MsgContacts)
	w.robot.LeftClick(MsgContacts)

	// Go to More
	w.robot.MoveTo(w.config.More.X+w.config.More.MarginX, w.config.More.Y+w.config.More.MarginY, MsgMore)
	w.robot.LeftClick(MsgLastContact)
	w.robot.RightClick(MsgLastContact)

	// Press Esc to dismiss any active menu/dialog and set focus
	w.robot.PressEsc(MsgLastContact)
	w.robot.LeftClick(MsgLastContact)

	// Capture small image for end-of-list check
	if w.config.CheckEnd {
		x, y := robotgo.Location()
		w.lastContactImg = w.vision.CaptureContact(x, y)
		if w.lastContactImg != nil {
			bestMatch, bestPoints := w.vision.FindContactPosition(w.lastContactImg)
			if bestMatch > 0.9 {
				w.lastContactPos = bestPoints
				if w.config.Debug {
					i18n.P.Printf("Last contact pos: %v [match:%0.3f]\n", w.lastContactPos, bestMatch)
				}
			} else {
				i18n.P.Printf("Last contact not found: [match:%0.3f]\n", bestMatch)
			}
		}
	}
}

func (w *WeChat) GotoNextContact() bool {
	end := false
	if w.messagesSent <= 0 {
		return end
	}

	w.robot.PressDown(MsgLastContact)

	if w.config.CheckEnd {
		if !w.vision.CheckContactChanged(w.lastContactImg, w.lastContactPos) {
			i18n.P.Printf("Contact did not change (End of list)\n")
			end = true
		}
	}
	return end
}

func (w *WeChat) WaitForMessageBar() {
	if w.config.Debug {
		i18n.P.Printf("Waiting for messages bar to appear...\n")
	}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if w.vision.CheckMessagesBar() {
			if w.config.Debug {
				i18n.P.Printf("Messages bar detected.\n")
			}
			return
		}

		if w.config.Debug {
			i18n.P.Printf("Messages bar not detected yet, retrying... (%d/%d)\n", i+1, maxRetries)
		}
		w.robot.Wait("WaitForMessageBar Interval")
	}

	if w.config.Debug {
		i18n.P.Printf("Warning: Messages bar did not appear after %d retries.\n", maxRetries)
	}
}

func (w *WeChat) SendMessage() {
	// Go to More
	w.robot.MoveTo(w.config.More.X+w.config.More.MarginX, w.config.More.Y+w.config.More.MarginY, MsgMore)
	w.robot.RightClick(MsgLastContact)

	// Go to Messages Popup
	w.robot.MoveRelative(w.config.MessagesPopup.MarginX, w.config.MessagesPopup.MarginY, MsgMessagesPopup)
	w.robot.LeftClick(MsgMessagesPopup)

	// Wait for the messages input area to appear dynamically
	w.WaitForMessageBar()

	// Go to Messages Input area
	w.robot.MoveTo(w.config.MessagesBar.X+w.config.MessagesBar.MarginX, w.config.MessagesBar.Y+w.config.MessagesBar.MarginY, MsgMessagesInput)
	w.robot.LeftClick(MsgMessagesInput)

	// Paste and Send
	w.robot.PasteClipboard(MsgMessagesInput)
	if w.config.AutoSend {
		w.robot.PressEnter(MsgMessagesInput)
	}
}

func (w *WeChat) ClearMessage() {
	// Go to More
	w.robot.MoveTo(w.config.More.X+w.config.More.MarginX, w.config.More.Y+w.config.More.MarginY, MsgMore)
	w.robot.RightClick(MsgLastContact)

	// Go to Messages
	w.robot.MoveRelative(w.config.MessagesPopup.MarginX, w.config.MessagesPopup.MarginY, MsgMessagesPopup)
	w.robot.LeftClick(MsgMessagesPopup)

	// Wait for the messages input area to appear dynamically
	w.WaitForMessageBar()

	// Go to Messages Input area
	w.robot.MoveTo(w.config.MessagesBar.X+w.config.MessagesBar.MarginX, w.config.MessagesBar.Y+w.config.MessagesBar.MarginY, MsgMessagesInput)
	w.robot.LeftClick(MsgMessagesInput)

	// Clear input
	w.robot.PressCtrlA(MsgMessagesInput)
	w.robot.PressDelete(MsgMessagesInput)
}

func (w *WeChat) IncrementMessagesSent() {
	w.messagesSent++
}

func (w *WeChat) GetMessagesSent() int {
	return w.messagesSent
}
