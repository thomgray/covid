package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thomgray/covid/domain"
	"github.com/thomgray/egg"
	"github.com/thomgray/egg/eggc"
)

type SummaryView struct {
	*egg.View
	border         *eggc.BorderView
	countrySummary *domain.CountrySummary
	globalSummary  *domain.GlobalSummary
}

func MakeSummaryView() *SummaryView {
	sv := SummaryView{}
	sv.View = egg.MakeView()
	sv.border = eggc.MakeBorderView()
	sv.border.AddSubView(sv.View)
	sv.border.SetTitle("Summary")

	sv.OnDraw(sv.draw)
	return &sv
}

func (sv *SummaryView) SetVisible(vis bool) {
	sv.border.SetVisible(vis)
}

func (sv *SummaryView) SetBounds(bs egg.Bounds) {
	sv.border.SetBounds(bs)
	bs.X = 1
	bs.Y = 1
	bs.Width -= 2
	bs.Height -= 2
	sv.View.SetBounds(bs)
}

func (sv *SummaryView) GetView() *egg.View {
	return sv.border.GetView()
}

func (sv *SummaryView) draw(c egg.Canvas) {
	x := 5
	y := 2
	if sv.globalSummary != nil {
		c.DrawString("Global", x, y, egg.ColorMagenta, c.Background, egg.AttrBold|egg.AttrUnderline)
		y += 2
		c.DrawString2(fmt.Sprintf("Total Confirmed:       %s", withThousandSeparators(sv.globalSummary.TotalConfirmed)), x, y)
		y++
		c.DrawString2(fmt.Sprintf("New Confirmed:         %s", withThousandSeparators(sv.globalSummary.NewConfirmed)), x, y)
		y++

		c.DrawString2(fmt.Sprintf("Total Deaths:          %s", withThousandSeparators(sv.globalSummary.TotalDeaths)), x, y)
		y++
		c.DrawString2(fmt.Sprintf("New Deaths:            %s", withThousandSeparators(sv.globalSummary.NewDeaths)), x, y)
		y++

		c.DrawString2(fmt.Sprintf("Total Recovered:       %s", withThousandSeparators(sv.globalSummary.TotalRecovered)), x, y)
		y++
		c.DrawString2(fmt.Sprintf("New Recovered:         %s", withThousandSeparators(sv.globalSummary.NewRecovered)), x, y)
		y++
	}

	y = 15

	if sv.countrySummary != nil {
		c.DrawString(fmt.Sprintf("For %s", sv.countrySummary.Country), x, y, egg.ColorRed, c.Background, egg.AttrUnderline|egg.AttrBold)
		y += 2

		c.DrawString2(fmt.Sprintf("Total Confirmed:       %s", withThousandSeparators(sv.countrySummary.TotalConfirmed)), x, y)
		y++
		c.DrawString2(fmt.Sprintf("New Confirmed:         %s", withThousandSeparators(sv.countrySummary.NewConfirmed)), x, y)
		y++

		c.DrawString2(fmt.Sprintf("Total Deaths:          %s", withThousandSeparators(sv.countrySummary.TotalDeaths)), x, y)
		y++
		c.DrawString2(fmt.Sprintf("New Deaths:            %s", withThousandSeparators(sv.countrySummary.NewDeaths)), x, y)
		y++

		c.DrawString2(fmt.Sprintf("Total Recovered:       %s", withThousandSeparators(sv.countrySummary.TotalRecovered)), x, y)
		y++
		c.DrawString2(fmt.Sprintf("New Recovered:         %s", withThousandSeparators(sv.countrySummary.NewRecovered)), x, y)
		y++
	}
}

func withThousandSeparators(i int) string {
	s := strconv.Itoa(i)

	pieces := make([]string, 0, (len(s)/3)+1)
	for len(s) > 3 {
		x := len(s) - 3
		pieces = append([]string{s[x:]}, pieces...)
		s = s[:x]
		i++
	}
	if len(s) > 0 {
		pieces = append([]string{s}, pieces...)
	}

	return strings.Join(pieces, ",")
}
