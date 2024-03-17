package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
	"log"
	"main/api/spotify"
	"main/ledos"
	"main/ledos/gpio"
	"main/uilib"
	"os"

	"github.com/joho/godotenv"
)

var (
	mainLog = log.New(os.Stdout, "MAIN: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		screamAndDie(err, mainLog)
	}
}

func main() {
	var text string
	var err error

	loadDotEnv()
	encoderRotation := make(chan int)
	reader := bufio.NewReader(os.Stdin)

	// DEBUG
	// spotify.Init()
	// spotify.TriggerAuth()
	var img = make([][]color.Color, 10)
	for i := 0; i < 10; i++ {
		img[i] = make([]color.Color, 10)

		for j := 0; j < 10; j++ {
			img[i][j] = color.White
		}
	}

	// customColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// ledos.DrawIsoscelesTriangle(image.Point{X: 32, Y: 8}, 5, 1, customColor)

	var c1 = uilib.CreateComponent(ledos.Canvas, 10, 10, 10, 10, img)
	var c2 = uilib.CreateComponent(ledos.Canvas, 30, 10, 10, 10, img)
	c1.Render(ledos.Canvas)
	c2.Render(ledos.Canvas)
	c1.SimpleHighlight()
	ledos.Render()
while:
	for {
		text, err = reader.ReadString('\n')
		screamAndDie(err, mainLog)

		switch text {
		case "q\n":
			break while

		case "m\n":
			oc := gpio.CreateMainController(ledos.Canvas)

			oc.Run()
			ledos.Render()
		case "s\n":
			spotify.Init()
		case "gpio\n":
			go gpio.TestGpio(encoderRotation)
		case "tr\n":
			customColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}
			ledos.DrawIsoscelesTriangle(image.Point{X: 32, Y: 8}, 5, 1, customColor)
		case "image\n":
			customColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}
			ledos.FillColor(customColor)
		case "sakura\n":
			ledos.Sakura()
		default:
			ledos.Dashboard()
		}
	}
}

func init() {
	flag.Parse()
}

func screamAndDie(err error, logger *log.Logger) {
	if err != nil {
		logger.Panic(err)
	}
}
