package apps

import (
	"fmt"
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

func Spotify(cnv *rgbmatrix.Canvas) *uilib.UI {
	canvas = cnv
	spotify := uilib.CreateUI(cnv)
	albumImage := albumComponent()
	nextTrack := prevTrackComponent()
	pauseTrack := pauseTrackComponent()
	previousTrack := nextTrackComponent()
	trackBanner := trackBannerComponent()

	spotify.AddComponent(albumImage)
	spotify.AddComponent(nextTrack)
	spotify.AddComponent(pauseTrack)
	spotify.AddComponent(previousTrack)
	spotify.AddComponent(trackBanner)
	return &spotify
}

func albumComponent() *uilib.Component {
	refreshFn := func(c *uilib.Component) {
		playbackState := spotify.GetPlaybackState()
		trackId := playbackState.Item.Id
		if trackId != c.State["trackId"] {
			c.State["trackId"] = trackId
		} else {
			imgPath := api.DownloadImage(playbackState.Item.Album.Images[0].Url)
			img := utils.LoadImage(imgPath)
			resizedImg := resize.Resize(30, 30, img, resize.Lanczos3)
			c.StaticFill = ledos.DrawImage(resizedImg)

		}
	}
	playbackState := spotify.GetPlaybackState()
	imgPath := api.DownloadImage(playbackState.Item.Album.Images[0].Url)
	img := utils.LoadImage(imgPath)
	resizedImg := resize.Resize(30, 30, img, resize.Lanczos3)

	fill := ledos.DrawImage(resizedImg)

	track := uilib.CreateComponent(canvas, 1, 1, 30, 30, fill)
	track.SetState("trackId", playbackState.Item.Id)
	track.SetRefreshRate(2000)
	track.SetRefreshRoutine(refreshFn)
	return track
}

func prevTrackComponent() *uilib.Component {

	fill := make([][]color.Color, 4)
	for i := 0; i < 4; i++ {
		fill[i] = make([]color.Color, 7)

		if i == 0 {

			for j := 0; j < 7; j++ {
				fill[i][j] = color.White
			}
		} else {
			for j := 0; j < 7; j++ {
				fill[i][j] = color.Black
			}
			for j := 0; j < i+1; j++ {
				fill[i][3+j] = color.White
				fill[i][3-j] = color.White
			}
		}

	}
	c := uilib.CreateComponent(canvas, 35, 5, 4, 7, fill)
	c.Actions.Click = func() { spotify.PreviousTrack() }
	return c
}

func pauseTrackComponent() *uilib.Component {

	fill := make([][]color.Color, 7)
	for i := 0; i < 7; i++ {
		fill[i] = make([]color.Color, 7)
		if i == 3 {
			for j := 0; j < 7; j++ {
				fill[i][j] = color.Black
			}
		} else {
			for j := 0; j < 7; j++ {
				fill[i][j] = color.White
			}
		}
	}
	c := uilib.CreateComponent(canvas, 44, 5, 7, 7, fill)
	c.Actions.Click = func() { spotify.ResumePlayback() }
	return c
}

func nextTrackComponent() *uilib.Component {
	fill := make([][]color.Color, 4)
	for i := 0; i < 4; i++ {
		fill[i] = make([]color.Color, 7)
		if i == 3 {
			for j := 0; j < 7; j++ {
				fill[i][j] = color.White
			}
		} else {
			for j := 0; j < 7; j++ {
				fill[i][j] = color.White
			}
			for j := 0; j < i; j++ {
				fill[i][j] = color.Black
				fill[i][6-j] = color.Black
			}
		}
	}
	c := uilib.CreateComponent(canvas, 56, 5, 4, 7, fill)
	c.Actions.Click = func() { spotify.NextTrack() }
	return c
}

func trackBannerComponent() *uilib.Component {
	refreshFn := func(c *uilib.Component) {
		trackName := spotify.GetTrackName()
		if trackName != c.State["trackName"] {

			newImage := ledos.WriteText(trackName, color.White)
			c.SetState("trackName", trackName)
			c.SetState("baseImage", newImage)
			c.SetState("imageLen", len(newImage))
			c.SetState("scrollPosition", 0)
			c.StaticFill = newImage
		} else {
			pos, posOk := c.State["scrollPosition"].(int)
			imgLen, lenOk := c.State["imageLen"].(int)

			if posOk && lenOk {
				if pos == imgLen-1 {
					c.SetState("scrollPosition", 0)
				} else {
					c.SetState("scrollPosition", pos+1)
				}
			} else {
				fmt.Println("Failed to parse State.scrollPosition : ", c.State["scrollPosition"])
			}

			img, ok := c.State["baseImage"].([][]color.Color)
			if ok {
				var scrolledImg [][]color.Color
				fmt.Println("len(img):", len(img))
				fmt.Println("Pos, len(img):", pos, len(img))
				if pos+26 /*viewport size*/ > imgLen-1 {

					firstSegment, secondSegment := img[pos:], img[:(pos+26)-imgLen]
					scrolledImg = append(firstSegment, secondSegment...)
				} else {
					scrolledImg = img[pos : pos+26]
				}

				fmt.Println("len(scrolledImg):", len(scrolledImg))
				c.StaticFill = scrolledImg
			} else {
				fmt.Println("Failed to parse State.baseImage : ", c.State["baseImage"])
			}
		}
	}
	trackName := spotify.GetTrackName()
	fill := ledos.WriteText(trackName, color.White)

	banner := uilib.CreateComponent(canvas, 35, 20, 26, 5, fill)
	banner.SetRefreshRate(200)
	banner.SetState("trackName", trackName)
	banner.SetState("baseImage", fill)
	banner.SetState("imageLen", len(fill))
	banner.SetState("scrollPosition", 0)
	banner.SetRefreshRoutine(refreshFn)
	return banner
}
