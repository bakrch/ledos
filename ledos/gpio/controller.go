package gpio

import (
	"main/ledos/apps"
	"main/uilib"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

type OsController struct {
	canvas            *rgbmatrix.Canvas
	apps              []*uilib.UI
	SelectedApp       int
	componentsFocused bool
}

func CreateMainController(cnv *rgbmatrix.Canvas) OsController {
	var oc OsController
	oc.componentsFocused = false
	oc.canvas = cnv
	oc.apps = make([]*uilib.UI, 0)
	spotifyApp := apps.Spotify(oc.canvas)

	oc.apps = append(oc.apps, spotifyApp)
	oc.SelectedApp = 0
	return oc
}

func (oc *OsController) ExecuteComponentAction(input uilib.Action) {
	app := oc.CurrentApp()
	idx := app.ActiveComponent
	app.Components[idx].Actions[input]()
}

func (oc *OsController) CurrentApp() *uilib.UI {
	return oc.apps[oc.SelectedApp]
}

func (oc *OsController) SelectNextApp() {
	if oc.SelectedApp++; oc.SelectedApp == len(oc.apps) {
		oc.SelectedApp = 0
	}
}

func (oc *OsController) FocusComponent(index int) {
	// oc.CurrentApp().Components[index].SetActive()
}
func (oc *OsController) FocusComponents() {
	oc.componentsFocused = true
}

func (oc *OsController) FocusUI() {
	oc.componentsFocused = false
}

func (oc *OsController) Run() {
	oc.apps[oc.SelectedApp].Render()
}
