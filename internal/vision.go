package internal

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/activebook/auto-wechat/internal/i18n"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

// Vision handles computer vision tasks (OpenCV matching, screenshoting)
type Vision struct {
	robot    *Robot
	appTitle string
	debug    bool
}

func NewVision(robot *Robot, appTitle string, debug bool) *Vision {
	return &Vision{
		robot:    robot,
		appTitle: appTitle,
		debug:    debug,
	}
}

// In order to compare images in memory
// We must convert it to correct RGBA format
func (v *Vision) normalizeImage(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// Capture the app image
// If app is not found, capture the screen
func (v *Vision) CaptureAppImage() image.Image {
	foundApp := true
	ax, ay, aw, ah, err := v.robot.GetAppWindowBounds(v.appTitle)
	if err != nil {
		i18n.P.Printf("get app window bounds error: %v\n", err)
		foundApp = false
	}
	v.robot.MoveTo(-100, -100, MsgOutOfScreen) // hide cursor

	var captured image.Image
	if foundApp {
		captured, err = robotgo.CaptureImg(ax, ay, aw, ah)
		if err != nil {
			i18n.P.Printf("captureAppImage[%d, %d, %d, %d]: capture failed: %v\n", ax, ay, aw, ah, err)
		}
	}

	if captured == nil {
		captured, err = robotgo.CaptureImg()
		if err != nil || captured == nil {
			i18n.P.Printf("capture screen error: %v\n", err)
			return nil
		}
	}

	img := v.normalizeImage(captured)
	return img
}

// Capture the contact image at (x, y)
// We need to capture the contact image to find the contact position
// The contact image is centered on x,y and with 100px width and 40px height
func (v *Vision) CaptureContact(x, y int) image.Image {
	v.robot.MoveTo(-100, -100, MsgOutOfScreen) // hide cursor
	captured, err := robotgo.CaptureImg(x-50, y-20, 100, 40)
	if err != nil || captured == nil {
		i18n.P.Printf("CaptureContact: capture failed: %v\n", err)
		return nil
	}
	return v.normalizeImage(captured)
}

// Find the position of the contact image
// lastContactImg: the last contact image
func (v *Vision) FindContactPosition(lastContactImg image.Image) (float32, image.Point) {
	appImg := v.CaptureAppImage()
	if appImg == nil {
		return 0, image.Point{}
	}
	_, bestMatch, _, bestPoints := gcv.FindImg(lastContactImg, appImg)
	return bestMatch, bestPoints
}

// Check whether the contact list has changed, returns true if changed
// logic: if the contact list has changed, the last contact position will change
// if the contact list has not changed, the last contact position will not change
func (v *Vision) CheckContactChanged(lastContactImg image.Image, lastContactPos image.Point) bool {
	if lastContactImg == nil {
		return true
	}

	appImg := v.CaptureAppImage()
	if appImg == nil {
		return true
	}

	_, bestMatch, _, bestPoints := gcv.FindImg(lastContactImg, appImg)
	isSame := bestMatch > 0.9

	if isSame && bestPoints == lastContactPos {
		return false // No change
	}
	return true // Changed
}

func (v *Vision) FindPoint(images []string) (Point, error) {
	screenshot, _ := robotgo.CaptureImg()

	for _, imgPath := range images {
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			if v.debug {
				fmt.Printf("Image file does not exist: %s\n", imgPath)
			}
			continue
		}

		template, _, err := robotgo.DecodeImg(imgPath)
		if err != nil {
			if v.debug {
				fmt.Printf("Error decoding image %s: %v\n", imgPath, err)
			}
			continue
		}

		_, bestMatch, _, bestPoint := gcv.FindImg(template, screenshot)

		if bestMatch > 0.8 {
			bounds := template.Bounds()
			return Point{X: bestPoint.X, Y: bestPoint.Y, W: bounds.Dx(), H: bounds.Dy()}, nil
		}
	}

	return Point{}, fmt.Errorf("image not found in vision")
}

// Check whether the messages bar is shown (after clicking on messages-send)
func (v *Vision) CheckMessagesBar() bool {
	appDir := GetAppDir()
	messagesbarImgs := []string{
		filepath.Join(appDir, "imgs/messages_bar_light.png"),
		filepath.Join(appDir, "imgs/messages_bar_dark.png"),
	}

	// Returns nil error if found
	_, err := v.FindPoint(messagesbarImgs)
	return err == nil
}

// Check whether the messages popup is shown (after right-clicking on contact)
func (v *Vision) CheckMessagesPopup() bool {
	appDir := GetAppDir()
	messagesPopupImgs := []string{
		filepath.Join(appDir, "imgs/messages_send_en_light.png"),
		filepath.Join(appDir, "imgs/messages_send_en_dark.png"),
		filepath.Join(appDir, "imgs/messages_send_zh_light.png"),
		filepath.Join(appDir, "imgs/messages_send_zh_dark.png"),
	}

	// Returns nil error if found
	_, err := v.FindPoint(messagesPopupImgs)
	return err == nil
}
