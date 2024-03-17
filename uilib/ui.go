package uilib

import (
	"fmt"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

var (
	canvas *rgbmatrix.Canvas
)

type UI struct {
	Components []*Component
	Active     bool
}

func CreateImage() {

}

func CreateUI(cnv *rgbmatrix.Canvas) UI {
	var ui UI
	ui.Active = true
	ui.Components = make([]*Component, 0)
	canvas = cnv
	return ui
}

func (ui *UI) AddComponent(c *Component) {
	fmt.Println("Component: ", c)
	ui.Components = append(ui.Components, c)
}

func (ui *UI) Render() {
	fmt.Println("UI: ", ui)
	for _, c := range ui.Components {
		if c != nil {
			c.Render(canvas)
		}
	}
}
