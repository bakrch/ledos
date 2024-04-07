package uilib

import (
	"time"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

var (
	canvas *rgbmatrix.Canvas
)

type UI struct {
	Components  []*Component
	Active      bool
	RefreshRate int
}

func CreateImage() {

}

func CreateUI(cnv *rgbmatrix.Canvas) UI {
	var ui UI
	ui.Active = true
	ui.RefreshRate = 100
	ui.Components = make([]*Component, 0)
	canvas = cnv
	return ui
}

func (ui *UI) AddComponent(c *Component) {
	ui.Components = append(ui.Components, c)
}

func (ui *UI) Render() {

	for _, c := range ui.Components {
		if c != nil {
			c.Render(canvas)
		}
	}
	done := make(chan bool)
	stop := make(chan bool)

	for _, c := range ui.Components {
		if c.RefreshRate != 0 {
			go func(comp *Component) {
				ticker := time.NewTicker(time.Duration(comp.RefreshRate) * time.Millisecond)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						comp.Refresh(comp)
						for _, c := range ui.Components {
							if c != nil {
								c.Render(canvas)
							}
						}
					case <-stop:
						done <- true
						return

					}
				}
			}(c)
		}
	}
}
