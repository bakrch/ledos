package apps

import (
	"image"
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
	albumImageComp := trackComponent()
	spotify.AddComponent(&albumImageComp)
	return &spotify
}

func trackComponent() uilib.Component {

	playbackState := spotify.GetPlaybackState()
	imgPath := api.DownloadImage(playbackState.Item.Album.Images[0].Url)
	img := utils.LoadImage(imgPath)
	resizedImg := resize.Resize(22, 22, img, resize.Lanczos3)
	fill := ledos.DrawImage(canvas, resizedImg, image.Point{21, 2})
	track := uilib.CreateComponent(canvas, 21, 2, 22, 20, fill)
	return track
}
