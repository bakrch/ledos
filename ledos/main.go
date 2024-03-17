package ledos

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

var (
	debugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	matrix   = initRgbMatrix()
	Canvas   = rgbmatrix.NewCanvas(matrix)
)

func Render() {
	Canvas.Render()
}

func Dashboard() {
	whiteColor := color.RGBA{R: 255, G: 0, B: 0, A: 255} // This is a lie
	writeText(Canvas, 4, 8, "GUTTEN TAG", whiteColor)
	writeText(Canvas, 4, 18, "GUTTEN TAG", whiteColor)
	writeText(Canvas, 4, 28, "GUTTEN TAG", whiteColor)
	Canvas.Render()
}

/*
*

			tip
			/|\
		   / | \
	   Sa /  |h \ Sb

Directions (of the tip) :

	  		 1
		 2		 3
			 4
*/
func DrawIsoscelesTriangle(tip image.Point, h int, direction int, color color.Color) {
	x := tip.X
	y := tip.Y
	for i := x; i < h+x; i++ {
		for j := 0; j < i-x; j++ {
			Canvas.Set(i, y+j, color)
			Canvas.Set(i, y-j, color)
		}
	}
}
func FillColor(color color.Color) {
	b := Canvas.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			Canvas.Set(x, y, color)
		}
	}
	Canvas.Render()
}

func Sakura() {
	customColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	// drawImage(Canvas, "./ledos/static/img/sakura-bg.png")
	writeText(Canvas, 4, 8, "GUTTEN TAG", customColor)
	Canvas.Render()
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println("Press Enter to exit...")
	// scanner.Scan()
}

func loadFont(path string) font.Face {
	fontData, err := os.ReadFile("./ledos/static/font/tiny.otf")
	fatal(err, debugLog)

	myFont, err := opentype.Parse(fontData)
	fatal(err, debugLog)
	face, err := opentype.NewFace(myFont, &opentype.FaceOptions{
		Size:    5,  // Set the font size
		DPI:     72, // Set the dots per inch
		Hinting: font.HintingFull,
	})
	fatal(err, debugLog)

	return face
}

func initRgbMatrix() rgbmatrix.Matrix {

	var (
		rows                     = flag.Int("led-rows", 32, "number of rows supported")
		cols                     = flag.Int("led-cols", 64, "number of columns supported")
		parallel                 = flag.Int("led-parallel", 1, "number of daisy-chained panels")
		chain                    = flag.Int("led-chain", 1, "number of displays daisy-chained")
		rgb_sequence             = flag.String("led-rgb-sequence", "RBG", "Order of colors on the matrix")
		brightness               = flag.Int("brightness", 100, "brightness (0-100)")
		hardware_mapping         = flag.String("led-gpio-mapping", "adafruit-hat-pwm", "Name of GPIO mapping used.")
		show_refresh             = flag.Bool("led-show-refresh", false, "Show refresh rate.")
		inverse_colors           = flag.Bool("led-inverse", false, "Switch if your matrix has inverse colors on.")
		disable_hardware_pulsing = flag.Bool("led-no-hardware-pulse", true, "Don't use hardware pin-pulse generation.")

		gpio_slowdown = flag.Int("led-slow-down-gpio", 2, "Slow down  GPIO")
	)
	config := &rgbmatrix.DefaultConfig
	config.Rows = *rows
	config.Cols = *cols
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness
	config.GPIOMapping = *hardware_mapping
	config.ShowRefreshRate = *show_refresh
	config.InverseColors = *inverse_colors
	config.DisableHardwarePulsing = *disable_hardware_pulsing
	config.RGBSequence = *rgb_sequence

	rtConfig := &rgbmatrix.DefaultRtConfig
	rtConfig.GPIOSlowdown = *gpio_slowdown
	matrix, err := rgbmatrix.NewRGBLedMatrix(config, rtConfig)

	fatal(err, debugLog)

	return matrix
}

func writeText(canvas draw.Image, x int, y int, text string, textColor color.Color) {

	f := loadFont("Hack-Regular.ttf")

	// Create a drawer
	d := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(textColor),
		Face: f,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}

	// Draw the text
	d.DrawString(text)
}

func DrawImage(canvas *rgbmatrix.Canvas, img image.Image, at image.Point) [][]color.Color {
	// Create a new RGBA image with the same bounds as the original image
	bounds := img.Bounds()
	imgArray := make([][]color.Color, bounds.Max.X-bounds.Min.X)
	// Convert each pixel from RGBA64 to RGBA
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		imgArray[x-bounds.Min.X] = make([]color.Color, bounds.Max.Y-bounds.Min.Y)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			imgArray[x-bounds.Min.X][y-bounds.Min.Y] = img.At(x, y)
		}
	}
	// draw.Draw(canvas, canvas.Bounds(), rgba, at, draw.Src)
	return imgArray
}

func fatal(err error, logger *log.Logger) {
	if err != nil {
		logger.Panic(err)
	}
}
