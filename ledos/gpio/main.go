package gpio

import (
	"fmt"
	"image"
	"image/color"
	"main/ledos"

	"github.com/stianeikeland/go-rpio/v4"
)

func TestGpio(direction <-chan int) {
	err := rpio.Open()
	defer rpio.Close()
	if err != nil {
		rpio.Close()
		panic(err)
	}
	dtPin := rpio.Pin(2)  // sda pin
	clkPin := rpio.Pin(3) //scl pin
	swPin := rpio.Pin(25)

	clkPin.High()
	dtPin.High()
	swPin.Low()

	var baseClk rpio.State = 1
	// var previousDt rpio.State = 1
	var (
		clk rpio.State
		dt  rpio.State
		sw  rpio.State
	)
	var counter = 0
	for {
		clk = clkPin.Read()
		dt = dtPin.Read()
		sw = swPin.Read()
		fmt.Println("dt, clk, sw", dt, clk, sw)
		if clk != baseClk {
			dt = dtPin.Read()
			if dtPin.Read() != clk {
				counter++
			} else {
				counter--
			}
			menuChoice := (counter) % 25
			if menuChoice < 0 {
				menuChoice = -1 * menuChoice
			}
			fmt.Println("dt ", dt, "clk ", clk, "previousClk", baseClk)
			customColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
			ledos.Dashboard()
			ledos.DrawIsoscelesTriangle(image.Point{X: 44, Y: menuChoice + 3}, 5, 1, customColor)

			// To prevent it from reading more than once per rotary encoder rotatio
			for clkPin.Read() != baseClk {
			}
			ledos.Render()
		}
	}
}
