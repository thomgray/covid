package main

import (
	"strings"

	"github.com/thomgray/egg/eggc"

	"github.com/thomgray/egg"
)

var app *egg.Application

var covidClient *CovidClient
var covidService *CovidServce
var dataController *DataController
var inputController *eggc.InputView
var inputLabel *eggc.LabelView
var output string
var countries []string

func main() {
	app = egg.InitOrPanic()
	defer app.Start()

	w, _ := egg.WindowSize()

	covidClient = (&CovidClient{}).Init()
	covidService = (&CovidServce{}).Init(covidClient)

	dataController = MakeDataController(app, covidService)
	inputController = eggc.MakeInputView()
	inputLabel = eggc.MakeLabelView()
	label := "Enter Country >"
	labelLen := len(label)
	inputLabel.SetLabel(label)
	inputLabel.SetForeground(egg.ColorMagenta)
	inputLabel.SetBounds(egg.MakeBounds(0, 0, labelLen, 1))

	countries = covidService.Countries()

	dataController.Init(countries)

	inputController.SetBounds(egg.MakeBounds(labelLen+1, 0, w, 1))

	app.AddViewController(inputController)
	app.SetFocusedView(inputController.GetView())
	app.AddViewController(inputLabel)

	inputController.Suggest(suggest)
	app.OnKeyEvent(keyHandler)

	app.OnResizeEvent(func(e *egg.ResizeEvent) {
		inputController.SetWidth(e.Width)
		dataController.Resize(e.Width, e.Height)
		app.ReDraw()
	})
}

func suggest(current string) (string, bool) {
	clean := strings.TrimSpace(current)
	if len(clean) == 0 {
		return "", false
	}

	for _, c := range countries {
		if strings.HasPrefix(c, clean) {
			remainder := strings.TrimPrefix(c, clean)
			return remainder, true
		}
	}
	return "", false
}

func keyHandler(e *egg.KeyEvent) {
	switch e.Key {
	case egg.KeyEsc:
		app.Stop()
	case egg.KeyEnter:
		e.StopPropagation = true
		input := strings.TrimSpace(inputController.GetTextContentString())
		works := dataController.LoadCountry(input)
		if works {
			inputController.SetTextContentString("")
			app.ReDraw()
		}
	default:
		switch e.Char {
		case '<', ',', '{', '[':
			e.StopPropagation = true
			dataController.SwitchPageLeft()
			app.ReDraw()
		case '>', '.', '}', ']':
			e.StopPropagation = true
			dataController.SwitchPageRight()
			app.ReDraw()
		}
	}
}
