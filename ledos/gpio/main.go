package gpio

import (
	"fmt"
	"main/ledos"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func InputController(direction <-chan int) {
	done := make(chan bool)
	stop := make(chan bool)
	GLOBAL_REFRESH_RATE := 500
	oc := CreateMainController(ledos.Canvas)

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
		//sw  rpio.State
	)
	var counter = 0
	for {
		clk = clkPin.Read()
		dt = dtPin.Read()
		//sw = swPin.Read()
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
			//fmt.Println("dt ", dt, "clk ", clk, "previousClk", baseClk)
			oc.CurrentApp().SelectNextComponent()
			fmt.Println("Current component: ", oc.CurrentApp().ActiveComponent)
			// To prevent it from reading more than once per rotary encoder rotatio
			for clkPin.Read() != baseClk {
			}
		}
	}
}
