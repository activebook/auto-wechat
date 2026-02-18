package test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/vcaesar/bitmap"
	"github.com/vcaesar/gcv"
	"github.com/vcaesar/imgo"
)

func Test() {
	test_mouse()
	// test_keyboard()
	// test_clipboard()
	// test_screen()
	// test_bitmap()
	// test_opencv()
	// test_opencv2()
	// test_hook()
	// test_win()
}

func test_mouse() {
	robotgo.MouseSleep = 300

	robotgo.Move(100, 100)
	fmt.Println(robotgo.Location())
	robotgo.Move(100, -200) // multi screen supported
	robotgo.MoveSmooth(120, -150)
	fmt.Println(robotgo.Location())

	robotgo.ScrollDir(10, "up")
	robotgo.ScrollDir(20, "right")

	robotgo.Scroll(0, -10)
	robotgo.Scroll(100, 0)

	robotgo.MilliSleep(100)
	robotgo.ScrollSmooth(-10, 6)
	// robotgo.ScrollRelative(10, -100)

	robotgo.Move(10, 20)
	robotgo.MoveRelative(0, -10)
	robotgo.DragSmooth(10, 10)

	robotgo.Click("wheelRight")
	robotgo.Click("left", true)
	robotgo.MoveSmooth(100, 200, 1.0, 10.0)

	robotgo.Toggle("left")
	robotgo.Toggle("left", "up")
}

func test_keyboard() {

	robotgo.KeyTap(robotgo.Tab, robotgo.Cmd)
	robotgo.Sleep(1)

	robotgo.Type("Hello World")
	robotgo.Type("だんしゃり", 0, 1)
	// robotgo.Type("テストする")

	robotgo.Type("Hi, Seattle space needle, Golden gate bridge, One world trade center.")
	robotgo.Type("Hi galaxy, hi stars, hi MT.Rainier, hi sea. こんにちは世界.")
	robotgo.Sleep(1)

	// ustr := uint32(robotgo.CharCodeAt("Test", 0))
	// robotgo.UnicodeType(ustr)

	robotgo.KeySleep = 100
	robotgo.KeyTap(robotgo.Enter)
	// // robotgo.Type("en")
	// robotgo.KeyTap("i", "alt", "cmd")

	// arr := []string{"alt", "cmd"}
	// robotgo.KeyTap("i", arr)

	// robotgo.MilliSleep(100)
	// robotgo.KeyToggle("a")
	// robotgo.KeyToggle("a", "up")
}

func test_clipboard() {
	// robotgo.WriteAll("Test")
	text, err := robotgo.ReadAll()
	if err == nil {
		fmt.Println(text)
	}
}

func test_screen() {
	x, y := robotgo.Location()
	fmt.Println("pos: ", x, y)

	color := robotgo.GetPixelColor(100, 200)
	fmt.Println("color---- ", color)

	sx, sy := robotgo.GetScreenSize()
	fmt.Println("get screen size: ", sx, sy)

	bit := robotgo.CaptureScreen(10, 10, 30, 30)
	defer robotgo.FreeBitmap(bit)

	img := robotgo.ToImage(bit)
	imgo.Save("test.png", img)

	num := robotgo.DisplaysNum()
	for i := 0; i < num; i++ {
		robotgo.DisplayID = i
		img1, _ := robotgo.CaptureImg()
		path1 := "save_" + strconv.Itoa(i)
		robotgo.Save(img1, path1+".png")
		robotgo.SaveJpeg(img1, path1+".jpeg", 50)

		img2, _ := robotgo.CaptureImg(10, 10, 20, 20)
		robotgo.Save(img2, "test_"+strconv.Itoa(i)+".png")

		x, y, w, h := robotgo.GetDisplayBounds(i)
		img3, err := robotgo.CaptureImg(x, y, w, h)
		fmt.Println("Capture error: ", err)
		robotgo.Save(img3, path1+"_1.png")
	}
}

func test_bitmap() {
	// bit := robotgo.CaptureScreen(10, 20, 30, 40)
	// // use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	// defer robotgo.FreeBitmap(bit)

	// fmt.Println("bitmap...", bit)
	// img := robotgo.ToImage(bit)
	// // robotgo.SavePng(img, "test_1.png")
	// robotgo.Save(img, "test_1.png")

	// bit2 := robotgo.ToCBitmap(robotgo.ImgToBitmap(img))

	template := bitmap.Open("contacts.png")
	fmt.Println("template: ", template) // if nil or empty, file not loading

	fx, fy := bitmap.Find(template)
	fmt.Println("FindBitmap------ ", fx, fy)
	robotgo.Move(fx, fy)
	robotgo.MilliSleep(100)
	robotgo.Click()

	arr := bitmap.FindAll(template)
	fmt.Println("Find all bitmap: ", arr)

	bitmap.Save(template, "test.png")
}

func test_opencv() {
	name := "test.png"
	name1 := "test_001.png"
	robotgo.SaveCapture(name1, 10, 10, 30, 30)
	robotgo.SaveCapture(name)

	fmt.Print("gcv find image: ")
	fmt.Println(gcv.FindImgFile(name1, name))
	fmt.Println(gcv.FindAllImgFile(name1, name))

	bit := bitmap.Open(name1)
	defer robotgo.FreeBitmap(bit)
	fmt.Print("find bitmap: ")
	fmt.Println(bitmap.Find(bit))

	// bit0 := robotgo.CaptureScreen()
	// img := robotgo.ToImage(bit0)
	// bit1 := robotgo.CaptureScreen(10, 10, 30, 30)
	// img1 := robotgo.ToImage(bit1)
	// defer robotgo.FreeBitmapArr(bit0, bit1)
	img, _ := robotgo.CaptureImg()
	img1, _ := robotgo.CaptureImg(10, 10, 30, 30)

	fmt.Print("gcv find image: ")
	fmt.Println(gcv.FindImg(img1, img))
	fmt.Println()

	res := gcv.FindAllImg(img1, img)
	fmt.Println(res[0].TopLeft.Y, res[0].Rects.TopLeft.X, res)
	x, y := res[0].TopLeft.X, res[0].TopLeft.Y
	robotgo.Move(x, y-rand.Intn(5))
	robotgo.MilliSleep(100)
	robotgo.Click()

	res = gcv.FindAll(img1, img) // use find template and sift
	fmt.Println("find all: ", res)
	res1 := gcv.Find(img1, img)
	fmt.Println("find: ", res1)

	img2, _, _ := robotgo.DecodeImg("test_001.png")
	x, y = gcv.FindX(img2, img)
	fmt.Println(x, y)
}

// 1.0  = perfect match
// 0.8+ = very good match
// 0.5+ = moderate match
// 0.42 = pretty weak match ← your result
// 0.0  = no correlation
// negative = inverse match
func test_opencv2() {
	robotgo.ActiveName("WeChat")
	img, _ := robotgo.CaptureImg()
	img1, format, _ := robotgo.DecodeImg("imgs/last_one.png")
	fmt.Println("format: ", format)
	fmt.Println("gcv find image: ")
	least_match, best_match, least_points, best_points := gcv.FindImg(img1, img)
	if best_match > 0.8 {
		fmt.Printf("Found match: %v\n", best_match)
		fmt.Println(best_points)
	} else if least_match < -0.8 {
		fmt.Printf("Found inverse match: %v\n", least_match)
		fmt.Println(least_points)
	} else {
		fmt.Printf("No match!\n")
	}
	fmt.Println()
}

func test_hook() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		hook.End()
	})

	fmt.Println("--- Please press w---")
	hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("w")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func test_win() {
	// fpid, err := robotgo.FindIds("WeChat")
	// if err == nil {
	// 	fmt.Println("pids... ", fpid)

	// 	if len(fpid) > 0 {
	// 		// robotgo.Type("Hi galaxy!", fpid[0])
	// 		// robotgo.KeyTap("a", fpid[0], "cmd")

	// 		// robotgo.KeyToggle("a", fpid[0])
	// 		// robotgo.KeyToggle("a", fpid[0], "up")

	// 		robotgo.ActivePid(fpid[0])

	// 		// robotgo.Kill(fpid[0])
	// 	}
	// }

	err := robotgo.ActiveName("Chrome")
	if err != nil {
		fmt.Println("active chrome error: ", err)
	}

	err = robotgo.ActiveName("Google Chrome")
	if err != nil {
		fmt.Println("active google chrome error: ", err)
	}

	err = robotgo.ActiveName("WeChat")
	if err != nil {
		fmt.Println("active wechat error: ", err)
	}
	fpid, err := robotgo.FindIds("WeChat")
	if err != nil {
		fmt.Println("find wechat error: ", err)
	}
	fmt.Println("pids... ", fpid)

	// isExist, err := robotgo.PidExists(100)
	// if err == nil && isExist {
	// 	fmt.Println("pid exists is", isExist)

	// 	robotgo.Kill(100)
	// }

	// abool := robotgo.Alert("test", "robotgo")
	// if abool {
	// 	fmt.Println("ok@@@ ", "ok")
	// }

	// title := robotgo.GetTitle()
	// fmt.Println("title@@@ ", title)
}
