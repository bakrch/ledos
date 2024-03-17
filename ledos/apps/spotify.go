package apps

import (
	"fmt"
	"image"
	"image/color"
	"main/api"
	"main/api/spotify"
	"main/ledos"
	"main/uilib"
	"main/uilib/utils"

	resize "github.com/nfnt/resize"
	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

var (
	canvas *rgbmatrix.Canvas
)

// spotify.Init()
// spotify.TriggerAuth()
func Spotify(cnv *rgbmatrix.Canvas) *uilib.UI {
	canvas = cnv
	spotify := uilib.CreateUI(cnv)
	albumImage := albumComponent()
	pauseTrack := pauseTrackComponent()
	nextTrack := nextTrackComponent()
	previousTrack := previousTrackComponent()

	spotify.AddComponent(&albumImage)
	spotify.AddComponent(&pauseTrack)
	spotify.AddComponent(&nextTrack)
	spotify.AddComponent(&previousTrack)
	return &spotify
}

func albumComponent() uilib.Component {

	playbackState := spotify.GetPlaybackState()
	imgPath := api.DownloadImage(playbackState.Item.Album.Images[0].Url)
	img := utils.LoadImage(imgPath)
	resizedImg := resize.Resize(32, 32, img, resize.Lanczos3)

	fill := ledos.DrawImage(canvas, resizedImg, image.Point{21, 2})

	track := uilib.CreateComponent(canvas, 0, 0, 32, 32, fill)
	return track
}

func pauseTrackComponent() uilib.Component {

	fill := make([][]color.Color, 4)
	for i := 0; i < 4; i++ {
		fill[i] = make([]color.Color, 7)
		for j := 0; j < 7; j++ {
			fill[i][j] = color.Black
		}
		for j := 0; j < i; j++ {
			fill[i][3+j] = color.White
			fill[i][3-j] = color.White
		}
	}
	fmt.Println("fill", fill)
	return uilib.CreateComponent(canvas, 33, 5, 4, 7, fill)
}

func nextTrackComponent() uilib.Component {

	fill := make([][]color.Color, 4)
	for i := 0; i < 4; i++ {
		fill[i] = make([]color.Color, 7)
		for j := 0; j < 7; j++ {
			fill[i][j] = color.Black
		}
		for j := 0; j < i; j++ {
			fill[i][3+j] = color.White
			fill[i][3-j] = color.White
		}
	}
	fmt.Println("fill", fill)
	return uilib.CreateComponent(canvas, 43, 5, 4, 7, fill)
}

func previousTrackComponent() uilib.Component {
	fill := make([][]color.Color, 4)
	for i := 0; i < 4; i++ {
		fill[i] = make([]color.Color, 7)
		for j := 0; j < 7; j++ {
			fill[i][j] = color.Black
		}
		for j := 0; j < i; j++ {
			fill[i][3+j] = color.White
			fill[i][3-j] = color.White
		}
	}
	fmt.Println("fill", fill)
	return uilib.CreateComponent(canvas, 53, 5, 4, 7, fill)
}
