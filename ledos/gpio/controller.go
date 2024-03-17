package gpio

import (
	"main/ledos/apps"
	"main/uilib"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

type OsController struct {
	canvas      *rgbmatrix.Canvas
	apps        []*uilib.UI
	selectedApp *uilib.UI
}

func CreateMainController(cnv *rgbmatrix.Canvas) OsController {
	var oc OsController
	oc.canvas = cnv
	oc.apps = make([]*uilib.UI, 0)
	spotifyApp := apps.Spotify(oc.canvas)

	oc.apps = append(oc.apps, spotifyApp)
	oc.selectedApp = oc.apps[0]
	return oc
}

func (oc *OsController) Run() {
	oc.selectedApp.Render()
}
