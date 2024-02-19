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
	spotify.Init()
	spotify.TriggerAuth()
while:
	for {
		text, err = reader.ReadString('\n')
		screamAndDie(err, mainLog)

		switch text {
		case "q\n":
			break while

		case "t\n":
			spotify.DoStuff()
		case "s\n":
			spotify.Init()
			spotify.DoStuff()
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
