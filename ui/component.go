package ui

import (
	"fmt"
	"image"
	"image/color"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

type Component struct {
	Placement   image.Rectangle
	StaticFill  [][]color.Color
	DynamicFill func() [][]color.Color
	Selected    bool
}

func (c Component) AmazingFunction() int {
	return c.Placement.Max.X - c.Placement.Min.X
}

func CreateComponent(x int, y int, width int, height int, fill [][]color.Color) Component {
	var c Component
	c.Placement.Min.X = x
	c.Placement.Min.Y = y
	c.Placement.Max.X = x + width
	c.Placement.Max.Y = y + height
	c.StaticFill = make([][]color.Color, width)
	for i := 0; i < width; i++ {
		c.StaticFill[i] = make([]color.Color, height)

		for j := 0; j < height; j++ {
			c.StaticFill[i][j] = fill[i][j]
		}
	}
	return c
}

func interpolateColor(c1, c2 color.RGBA, t float64) color.RGBA {
	r := uint8(float64(c1.R) + t*(float64(c2.R)-float64(c1.R)))
	g := uint8(float64(c1.G) + t*(float64(c2.G)-float64(c1.G)))
	b := uint8(float64(c1.B) + t*(float64(c2.B)-float64(c1.B)))
	a := uint8(float64(c1.A) + t*(float64(c2.A)-float64(c1.A)))
	return color.RGBA{r, g, b, a}
}

func (c Component) Highlight(canvas *rgbmatrix.Canvas) {
	rainbow := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{255, 127, 0, 255}, // Orange
		{255, 255, 0, 255}, // Yellow
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{75, 0, 130, 255},  // Indigo
		{148, 0, 211, 255}, // Violet
	}
	contourLen := c.Placement.Dx() + c.Placement.Dy()
	colors := make([]color.RGBA, contourLen)
	bandwidth := contourLen / (len(rainbow) - 1)

	for i := 0; i < bandwidth; i++ {
		t := float64(i%bandwidth) / float64(bandwidth-1)
		clr := interpolateColor(rainbow[i], rainbow[i+1], t)
		colors = append(colors, clr)
	}

	for x := c.Placement.Min.X; x < c.Placement.Max.X; x++ {
		canvas.Set(x, 0, colors[x])
		canvas.Set(x, c.Placement.Max.Y, colors[x])
	}

	for y := c.Placement.Min.Y; y < c.Placement.Max.Y; y++ {
		canvas.Set(0, y, colors[y])
		canvas.Set(c.Placement.Max.X, y, colors[y])
	}
}

// Only draws to the canvas, doesn't actually render the component until canvas.Render() is called
func (c Component) Render(canvas *rgbmatrix.Canvas) {
	i := 0
	for y := c.Placement.Min.Y; y < c.Placement.Max.Y; y++ {
		j := 0
		for x := c.Placement.Min.X; x < c.Placement.Max.X; x++ {
			fmt.Println("x, y", i, i)
			fmt.Println("canvas", canvas)
			canvas.Set(x, y, c.At(i, j))
			j++
		}
		i++
	}
}

func (c Component) At(x int, y int) color.Color {
	return c.StaticFill[x][y]
}
