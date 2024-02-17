package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
	"log"
	"main/ledos"
	"os"
)

var (
	mainLog = log.New(os.Stdout, "MAIN: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	var text string
	var err error

	reader := bufio.NewReader(os.Stdin)

while:
	for {
		text, err = reader.ReadString('\n')
		fatal(err, mainLog)
		switch text {
		case "q\n":
			break while

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

func fatal(err error, logger *log.Logger) {
	if err != nil {
		logger.Panic(err)
	}
}
