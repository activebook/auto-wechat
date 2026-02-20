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

func (v *Vision) normalizeImage(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

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

func (v *Vision) CaptureContact(x, y int) image.Image {
	v.robot.MoveTo(-100, -100, MsgOutOfScreen) // hide cursor
	captured, err := robotgo.CaptureImg(x-50, y-20, 100, 40)
	if err != nil || captured == nil {
		i18n.P.Printf("CaptureContact: capture failed: %v\n", err)
		return nil
	}
	return v.normalizeImage(captured)
}

func (v *Vision) FindContactPosition(lastContactImg image.Image) (float32, image.Point) {
	appImg := v.CaptureAppImage()
	if appImg == nil {
		return 0, image.Point{}
	}
	_, bestMatch, _, bestPoints := gcv.FindImg(lastContactImg, appImg)
	return bestMatch, bestPoints
}

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
			return Point{X: bestPoint.X, Y: bestPoint.Y}, nil
		}
	}

	return Point{}, fmt.Errorf("image not found in vision")
}

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
