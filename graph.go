package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/thomgray/egg/eggc"

	"github.com/thomgray/covid/domain"
	"github.com/thomgray/egg"
)

type GraphViewController struct {
	*egg.View
	border    *eggc.BorderView
	values    []domain.ValueByDay
	errorTile string
	isError   bool
}

const (
	oneEighth    rune = '▁'
	twoEighths        = '▂'
	threeEighths      = '▃'
	fourEighths       = '▄'
	fiveEighths       = '▅'
	sixEighths        = '▆'
	sevenEighths      = '▇'
	eightEighths      = '█'
)

func (cd *GraphViewController) Init() *GraphViewController {
	cd.View = egg.MakeView()
	cd.border = eggc.MakeBorderView()
	cd.border.AddSubView(cd.View)

	cd.View.OnDraw(cd.draw)
	return cd
}

// SetBounds - override from View
func (cd *GraphViewController) SetBounds(b egg.Bounds) {
	cd.border.SetBounds(b)
	b.Height -= 2
	b.X = 1
	b.Y = 1
	b.Width -= 2
	cd.View.SetBounds(b)
}

func (cd *GraphViewController) GetView() *egg.View {
	return cd.border.GetView()
}

func (cd *GraphViewController) SetData(values []domain.ValueByDay) {
	cd.values = values

}

func (cd *GraphViewController) SetErrorMessage(message string) {
	cd.errorTile = message
}

func (cd *GraphViewController) SetIsError(isError bool) {
	cd.isError = isError
}

func (cd *GraphViewController) SetVisible(vis bool) {
	cd.border.SetVisible(vis)
}

func (cd *GraphViewController) draw(c egg.Canvas) {
	if cd.isError {
		c.DrawString2(cd.errorTile, 0, 0)
	}
	if len(cd.values) == 0 {
		return
	}

	var highestVal int = 0
	vspace := c.Height - 2

	for _, v := range cd.values {
		if highestVal < v.Value {
			highestVal = v.Value
		}
	}
	yMax, increments := findYCeilWithIncrements(highestVal, vspace)
	incrementLabels := make([]string, len(increments))
	maxLabelWidth := 0
	labelSpacing := vspace / len(increments)
	for i, inc := range increments {
		label := formatYAxis(inc)
		incrementLabels[i] = label
		if runewidth.StringWidth(label) > maxLabelWidth {
			maxLabelWidth = runewidth.StringWidth(label)
		}
	}

	for i := vspace; i >= 0; i-- {
		c.DrawRune2('│', maxLabelWidth, i)
	}
	for i, label := range incrementLabels {
		c.DrawString2(label+"│̲", 0, vspace-((i+1)*labelSpacing))
	}

	var yscale float64 = float64(vspace*8) / float64(yMax)
	xoff := maxLabelWidth + 1
	graphW := c.Width - xoff
	dataSlice := cd.values
	sliceLen := len(dataSlice)
	if sliceLen > graphW {
		dataSlice = dataSlice[sliceLen-graphW:]
	}

	for x, v := range dataSlice {
		fg := egg.ColorBlue
		if x%2 == 0 {
			fg = egg.ColorGreen
		}
		barHeight := int(yscale * float64(v.Value))
		actualHeight := barHeight / 8
		remainder := barHeight % 8
		realx := x + xoff

		if x%10 == 0 {
			c.DrawString2(v.Day.Format("1/2"), realx, c.Height-1)
		}

		y := c.Height - 2
		for yi := 0; yi < actualHeight; yi++ {
			c.DrawRune(eightEighths, realx, y, fg, c.Background, c.Attribute)
			y--
		}
		if remainder > 0 {
			var char rune
			switch remainder {
			case 1:
				char = oneEighth
			case 2:
				char = twoEighths
			case 3:
				char = threeEighths
			case 4:
				char = fourEighths
			case 5:
				char = fiveEighths
			case 6:
				char = sixEighths
			case 7:
				char = sevenEighths
			default:
				char = ' '
			}
			c.DrawRune(char, realx, y, fg, c.Background, c.Attribute)
		}
	}
}

func findYCeil(maxVal int) int {
	ceil := 10
	for ceil < maxVal {
		ceil *= 10
	}
	// ceil is the lower v > maxVal power of 10, i.e: 10, 100, 1000, 10000 etc
	// but we want to return either
	tenth := ceil / 10

	for i := 0; i < 9; i++ {
		if maxVal < i*tenth {
			return i * tenth
		}
	}
	return ceil
}

func findYCeilWithIncrements(maxVal int, height int) (int, []int) {
	ceil := 10
	for ceil < maxVal {
		ceil *= 10
	}
	// ceil is the lower v > maxVal power of 10, i.e: 10, 100, 1000, 10000 etc
	// but we want to return either
	tenth := ceil / 10

	for i := 0; i < 9; i++ {
		if maxVal < i*tenth {
			return i * tenth, []int{i * tenth}
		}
	}
	return ceil, []int{ceil}
}

func formatYAxis(v int) string {
	if v < 1000 {
		return strconv.Itoa(v)
	} else if v < 500000 {
		formatted := fmt.Sprintf("%.1fk", float64(v)/1000.0)
		return strings.ReplaceAll(formatted, ".0k", "k")
	} else {
		formatted := fmt.Sprintf("%.1fM", float64(v)/1000000.0)
		return strings.ReplaceAll(formatted, ".0M", "M")
	}
}
