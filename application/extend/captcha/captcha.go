package captcha

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	dpi        = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	r          = rand.New(rand.NewSource(time.Now().UnixNano()))
	fontFamily = make([]string, 0)
)

const txtChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

const (
	//图片格式
	ImageFormatPng ImageFormat = iota
	ImageFormatJpeg
	ImageFormatGif
)
const (
	//验证码噪点强度
	CaptchaComplexLower CaptchaComplex = iota
	CaptchaComplexMedium
	CaptchaComplexHigh
)

type ImageFormat int

type CaptchaComplex int

type CaptchaImage struct {
	nrgba   *image.NRGBA
	width   int
	height  int
	Complex int
	Error   error
}

//go:embed fonts
var _fs embed.FS
var ioFS fs.FS = _fs

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_ = readFontFromFS(_fs)
}

/**
输出一个验证码图片，并返回图片的随机字符串用于session或redis存储
参数：w: 实现了write方法的接口
	width：宽度
	height：高度
	arg[0]:字符串个数，传为4个
	arg[1]:图片难辨度，:默认简单
*/
func NewCaptcha(w io.Writer, width, height int, arg ...int) (code string) {
	captchaImage := New(width, height, RandLightColor())

	if len(arg) > 1 {
		switch arg[1] {
		case 1:
			captchaImage.DrawNoise(CaptchaComplexLower)
			captchaImage.DrawTextNoise(CaptchaComplexLower)
			break
		case 2:
			captchaImage.DrawNoise(CaptchaComplexMedium)
			captchaImage.DrawTextNoise(CaptchaComplexMedium)
			break
		case 3:
			captchaImage.DrawNoise(CaptchaComplexHigh)
			captchaImage.DrawTextNoise(CaptchaComplexHigh)
			break
		default:
			captchaImage.DrawNoise(CaptchaComplexLower)
			captchaImage.DrawTextNoise(CaptchaComplexLower)
			break
		}
	} else {
		captchaImage.DrawNoise(CaptchaComplexLower)
		captchaImage.DrawTextNoise(CaptchaComplexLower)
	}

	var count int
	if len(arg) > 0 {
		count = arg[0]
	} else {
		count = 4
	}

	code = RandText(count)
	captchaImage.DrawText(code)
	captchaImage.DrawBorder(ColorToRGB(0x17A7A7A))
	captchaImage.DrawSineLine()
	captchaImage.SaveImage(w, ImageFormatPng)
	return code
}

func FileSystem(fs embed.FS) error {
	err := readFontFromFS(fs)
	if err != nil {
		return err
	}
	ioFS = fs
	return nil
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ReadFonts(dirPth string, suffix string) (err error) {
	fontFamily = fontFamily[:0]

	dir, err := ioutil.ReadDir(dirPth)

	if err != nil {
		return err
	}
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			fontFamily = append(fontFamily, filepath.Join(dirPth, fi.Name()))
		}
	}
	ioFS = os.DirFS(dirPth)
	return nil
}

//获取所及字体.
func RandFontFamily() (*truetype.Font, error) {
	fontFile := fontFamily[r.Intn(len(fontFamily))]
	_, path, _, _ := runtime.Caller(1)                       //获取挡圈目录
	fontFile2 := filepath.Join(filepath.Dir(path), fontFile) //链接目录到文件
	//fontBytes, err := ReadFile(fontFile)
	fontBytes, err := ReadFile(fontFile2)
	if err != nil {
		log.Printf("读取文件失败 -> %s - %+v\n", fontFile, err)
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("解析字体内容试标 -> %s - %+v\n", fontFile, err)
		return nil, err
	}
	return f, nil
}

//添加一个字体路径到字体库.
func AddFontFamily(fontPath ...string) {
	fontFamily = append(fontFamily, fontPath...)
}

//ReadFile 读取文件，优先从设置的文件系统中读取，失败后从内置文件系统读取，失败后从磁盘读取.
func ReadFile(name string) ([]byte, error) {
	if ioFS != nil {
		b, err := fs.ReadFile(ioFS, name)
		if err == nil {
			return b, nil
		}
	}
	return os.ReadFile(name)
}

func readFontFromFS(fs embed.FS) error {
	files, err := fs.ReadDir("fonts")
	if err != nil {
		log.Printf("解析字体文件失败 -> %+v", err)
		return err
	} else {
		fontFamily = fontFamily[:0]
		for _, fi := range files {
			if fi.IsDir() { // 忽略目录
				continue
			}
			if strings.HasSuffix(strings.ToLower(fi.Name()), ".ttf") { //匹配文件
				fontFamily = append(fontFamily, filepath.Join("fonts", fi.Name()))
			}
		}
	}
	return nil
}

//New 新建一个图片对象
func New(width int, height int, bgColor color.RGBA) *CaptchaImage {
	m := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)
	return &CaptchaImage{
		nrgba:  m,
		height: height,
		width:  width,
	}
}

//SaveImage 保存图片对象
func (captcha *CaptchaImage) SaveImage(w io.Writer, imageFormat ImageFormat) error {

	if imageFormat == ImageFormatPng {
		return png.Encode(w, captcha.nrgba)
	}
	if imageFormat == ImageFormatJpeg {
		return jpeg.Encode(w, captcha.nrgba, &jpeg.Options{Quality: 100})
	}
	if imageFormat == ImageFormatGif {
		return gif.Encode(w, captcha.nrgba, &gif.Options{NumColors: 256})
	}
	return errors.New("not supported image format")
}

//DrawHollowLine 添加一个较粗的空白直线
func (captcha *CaptchaImage) DrawHollowLine() *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	first := captcha.width / 20
	end := first * 19

	lineColor := color.RGBA{R: 245, G: 250, B: 251, A: 255}
	x1 := float64(r.Intn(first))
	x2 := float64(r.Intn(first) + end)
	multiple := float64(r.Intn(5)+3) / float64(5)
	if int(multiple*10)%3 == 0 {
		multiple = multiple * -1.0
	}
	w := captcha.height / 20
	for ; x1 < x2; x1++ {
		y := math.Sin(x1*math.Pi*multiple/float64(captcha.width)) * float64(captcha.height/3)
		if multiple < 0 {
			y = y + float64(captcha.height/2)
		}
		captcha.nrgba.Set(int(x1), int(y), lineColor)
		for i := 0; i <= w; i++ {
			captcha.nrgba.Set(int(x1), int(y)+i, lineColor)
		}
	}
	return captcha
}

//DrawSineLine 画一条曲线.
func (captcha *CaptchaImage) DrawSineLine() *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	px := 0
	var py float64 = 0

	//振幅
	a := r.Intn(captcha.height / 2)
	//Y轴方向偏移量
	b := Random(int64(-captcha.height/4), int64(captcha.height/4))
	//X轴方向偏移量
	f := Random(int64(-captcha.height/4), int64(captcha.height/4))
	// 周期
	var t float64 = 0
	if captcha.height > captcha.width/2 {
		t = Random(int64(captcha.width/2), int64(captcha.height))
	} else {
		t = Random(int64(captcha.height), int64(captcha.width/2))
	}
	w := (2 * math.Pi) / t

	// 曲线横坐标起始位置
	px1 := 0
	px2 := int(Random(int64(float64(captcha.width)*0.8), int64(captcha.width)))
	c := color.RGBA{R: uint8(r.Intn(150)), G: uint8(r.Intn(150)), B: uint8(r.Intn(150)), A: uint8(255)}
	for px = px1; px < px2; px++ {
		if w != 0 {
			py = float64(a)*math.Sin(w*float64(px)+f) + b + (float64(captcha.width) / float64(5))
			i := captcha.height / 5
			for i > 0 {
				captcha.nrgba.Set(px+i, int(py), c)
				i--
			}
		}
	}
	return captcha
}

type Point struct {
	X int
	Y int
}

//DrawLine 画一条直线.
func (captcha *CaptchaImage) DrawLine(num int) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	first := captcha.width / 10
	end := first * 9

	y := captcha.height / 3
	for i := 0; i < num; i++ {
		point1 := Point{X: r.Intn(first), Y: r.Intn(y)}
		point2 := Point{X: r.Intn(first) + end, Y: r.Intn(y)}
		if i%2 == 0 {
			point1.Y = r.Intn(y) + y*2
			point2.Y = r.Intn(y)
		} else {
			point1.Y = r.Intn(y) + y*(i%2)
			point2.Y = r.Intn(y) + y*2
		}
		captcha.drawBeeline(point1, point2, randDeepColor())
	}
	return captcha
}

//drawBeeline 画直线.
func (captcha *CaptchaImage) drawBeeline(point1 Point, point2 Point, lineColor color.RGBA) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	dx := math.Abs(float64(point1.X - point2.X))
	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y >= point2.Y {
		sy = -1
	}
	err := dx - dy
	//循环的画点直到到达结束坐标停止.
	for {
		captcha.nrgba.Set(point1.X, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X+1, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X-1, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X+2, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X-2, point1.Y, lineColor)
		if point1.X == point2.X && point1.Y == point2.Y {
			return captcha
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

//DrawBorder 画边框.
func (captcha *CaptchaImage) DrawBorder(borderColor color.RGBA) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	for x := 0; x < captcha.width; x++ {
		captcha.nrgba.Set(x, 0, borderColor)
		captcha.nrgba.Set(x, captcha.height-1, borderColor)
	}
	for y := 0; y < captcha.height; y++ {
		captcha.nrgba.Set(0, y, borderColor)
		captcha.nrgba.Set(captcha.width-1, y, borderColor)
	}
	return captcha
}

//DrawNoise 画噪点.
func (captcha *CaptchaImage) DrawNoise(complex CaptchaComplex) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	density := 18
	if complex == CaptchaComplexLower {
		density = 28
	} else if complex == CaptchaComplexMedium {
		density = 18
	} else if complex == CaptchaComplexHigh {
		density = 8
	}
	maxSize := (captcha.height * captcha.width) / density
	for i := 0; i < maxSize; i++ {
		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)
		captcha.nrgba.Set(rw, rh, randColor())
		size := r.Intn(maxSize)
		if size%3 == 0 {
			captcha.nrgba.Set(rw+1, rh+1, randColor())
		}
	}
	return captcha
}

//DrawTextNoise 画文字噪点.
func (captcha *CaptchaImage) DrawTextNoise(complex CaptchaComplex) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	density := 1500
	if complex == CaptchaComplexLower {
		density = 2000
	} else if complex == CaptchaComplexMedium {
		density = 1500
	} else if complex == CaptchaComplexHigh {
		density = 1000
	}

	maxSize := (captcha.height * captcha.width) / density
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetClip(captcha.nrgba.Bounds())
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)
	rawFontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(10))

	for i := 0; i < maxSize; i++ {
		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)
		text := RandText(1)
		fontSize := rawFontSize/2 + float64(r.Intn(5))
		c.SetSrc(image.NewUniform(RandLightColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			captcha.Error = err
			return captcha
		}
		c.SetFont(f)
		pt := freetype.Pt(rw, rh)
		_, err = c.DrawString(text, pt)
		if err != nil {
			captcha.Error = err
			return captcha
		}
	}
	return captcha
}

//DrawText 写字.
func (captcha *CaptchaImage) DrawText(text string) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetClip(captcha.nrgba.Bounds())
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)
	fontWidth := captcha.width / len(text)
	for i, s := range text {
		fontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(9))
		c.SetSrc(image.NewUniform(randDeepColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()
		if err != nil {
			captcha.Error = err
			return captcha
		}
		c.SetFont(f)
		x := (fontWidth)*i + (fontWidth)/int(fontSize)
		y := 5 + r.Intn(captcha.height/2) + int(fontSize/2)
		pt := freetype.Pt(x, y)
		_, err = c.DrawString(string(s), pt)
		if err != nil {
			captcha.Error = err
			return captcha
		}
		//pt.Y += c.PointToFixed(*size * *spacing)
		//pt.X += c.PointToFixed(*size);
	}
	return captcha

}

//生成指定大小的随机数.
func Random(min int64, max int64) float64 {
	if max <= min {
		panic(fmt.Sprintf("invalid range %d >= %d", max, min))
	}
	decimal := rand.Float64()
	if max <= 0 {
		return (float64(rand.Int63n((min*-1)-(max*-1))+(max*-1)) + decimal) * -1
	}
	if min < 0 && max > 0 {
		if rand.Int()%2 == 0 {
			return float64(rand.Int63n(max)) + decimal
		} else {
			return (float64(rand.Int63n(min*-1)) + decimal) * -1
		}
	}
	return float64(rand.Int63n(max-min)+min) + decimal
}

//随机生成深色系.
func randDeepColor() color.RGBA {
	randColor := randColor()
	increase := float64(30 + r.Intn(255))
	red := math.Abs(math.Min(float64(randColor.R)-increase, 255))
	green := math.Abs(math.Min(float64(randColor.G)-increase, 255))
	blue := math.Abs(math.Min(float64(randColor.B)-increase, 255))
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//随机生成浅色.
func RandLightColor() color.RGBA {
	red := r.Intn(55) + 200
	green := r.Intn(55) + 200
	blue := r.Intn(55) + 200
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//生成随机颜色.
func randColor() color.RGBA {
	red := r.Intn(255)
	green := r.Intn(255)
	blue := r.Intn(255)
	if (red + green) > 400 {
		blue = 0
	} else {
		blue = 400 - green - red
	}
	if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//生成随机字体.
func RandText(num int) string {
	textNum := len(txtChars)
	text := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		text = text + string(txtChars[r.Intn(textNum)])
	}
	return text
}

// 颜色代码转换为RGB
//input int
//output int red, green, blue.
func ColorToRGB(colorVal int) color.RGBA {
	red := colorVal >> 16
	green := (colorVal & 0x00FF00) >> 8
	blue := colorVal & 0x0000FF
	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(255),
	}
}

func (captcha *CaptchaImage) Save() (io io.Reader, err error) {
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, captcha.nrgba, nil)
	if err != nil {
		return
	}
	io = bytes.NewReader(buf.Bytes())
	return
}
