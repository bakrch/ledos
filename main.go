package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"main/api/spotify"
	"main/ledos"
	"main/ledos/gpio"
	"os"
	"os/exec"
	"time"

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
	var GLOBAL_REFRESH_RATE int = 400
	loadDotEnv()
	encoderRotation := make(chan int)
	reader := bufio.NewReader(os.Stdin)

	// DEBUG
	// spotify.Init()
	// spotify.TriggerAuth()
	ledos.Render()
	spotify.Init()
	cmd := exec.Command(os.Getenv("BROWSER"), []string{"--new-tab", "--url", "localhost:8080/login"}...)
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println("Error executing cmd", cmdErr)
	}

	done := make(chan bool)
	stop := make(chan bool)
loop:
	for {
		text, err = reader.ReadString('\n')
		screamAndDie(err, mainLog)

		switch text {
		case "q\n":
			break loop

		case "m\n":
			oc := gpio.CreateMainController(ledos.Canvas)

			oc.Run()

			go func() {
				ticker := time.NewTicker(time.Duration(GLOBAL_REFRESH_RATE) * time.Millisecond)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						ledos.Render()
					case <-stop:
						done <- true
						return
					}
				}
			}()
		case "s\n":
		case "gpio\n":
			go gpio.InputController(encoderRotation)
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
