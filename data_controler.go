package main

import (
	"github.com/thomgray/covid/domain"
	"github.com/thomgray/egg"
)

type DataController struct {
	covidService *CovidServce
	countries    []string
	summaries    domain.SummaryResponse
	app          *egg.Application

	pageSwitcher      *PageToggleViewController
	confirmedGraph    *GraphViewController
	confirmedNewGraph *GraphViewController
	deathsGraph       *GraphViewController
	deathsNewGraph    *GraphViewController
	recoveredGraph    *GraphViewController
	recoveredNewGraph *GraphViewController
	summaryView       *SummaryView
}

func MakeDataController(
	app *egg.Application,
	covidService *CovidServce,
) *DataController {
	dc := DataController{}
	dc.app = app
	dc.covidService = covidService
	dc.pageSwitcher = MakePageToggleViewController([]string{"Page 1", "Page 2"})
	return &dc
}

func (dc *DataController) SetPage(p int) {
	if p > 1 {
		return
	}
	dc.pageSwitcher.Page = p
	p0Visisble := p == 0
	p1Visisble := p == 1

	dc.confirmedGraph.SetVisible(p0Visisble)
	dc.confirmedNewGraph.SetVisible(p0Visisble)
	dc.deathsNewGraph.SetVisible(p1Visisble)
	dc.deathsGraph.SetVisible(p1Visisble)
	dc.recoveredGraph.SetVisible(p1Visisble)
	dc.recoveredNewGraph.SetVisible(p1Visisble)
	dc.summaryView.SetVisible(p0Visisble)
}

func (dc *DataController) SwitchPageLeft() {
	if dc.pageSwitcher.Page > 0 {
		dc.SetPage(dc.pageSwitcher.Page - 1)
	}
}

func (dc *DataController) SwitchPageRight() {
	if dc.pageSwitcher.Page < len(dc.pageSwitcher.PageNames)-1 {
		dc.SetPage(dc.pageSwitcher.Page + 1)
	}
}

func (dc *DataController) Init(countries []string) {
	dc.countries = countries
	w, h := egg.WindowSize()

	dc.pageSwitcher.SetBounds(egg.MakeBounds(0, 1, w, 1))

	dc.summaryView = MakeSummaryView()

	dc.confirmedGraph = (&GraphViewController{}).Init()
	dc.confirmedGraph.border.SetTitle("Total cases")

	dc.confirmedNewGraph = (&GraphViewController{}).Init()
	dc.confirmedNewGraph.border.SetTitle("Daily new cases")

	dc.deathsGraph = (&GraphViewController{}).Init()
	dc.deathsGraph.border.SetTitle("Total deaths")

	dc.deathsNewGraph = (&GraphViewController{}).Init()
	dc.deathsNewGraph.border.SetTitle("Daily deaths")

	dc.recoveredGraph = (&GraphViewController{}).Init()
	dc.recoveredGraph.border.SetTitle("Total recoveries")

	dc.recoveredNewGraph = (&GraphViewController{}).Init()
	dc.recoveredNewGraph.border.SetTitle("Daily recoveries")

	summaries, err := dc.covidService.Summaries()
	if err == nil {
		dc.summaries = summaries
		dc.summaryView.globalSummary = &dc.summaries.Global
	}

	app.AddViewController(dc.summaryView)
	app.AddViewController(dc.confirmedGraph)
	app.AddViewController(dc.confirmedNewGraph)
	app.AddViewController(dc.deathsGraph)
	app.AddViewController(dc.deathsNewGraph)
	app.AddViewController(dc.recoveredGraph)
	app.AddViewController(dc.recoveredNewGraph)

	app.AddViewController(dc.pageSwitcher)
	dc.Resize(w, h)

	dc.SetPage(0)
}

func (dc *DataController) countryExists(country string) bool {
	for _, c := range dc.countries {
		if c == country {
			return true
		}
	}
	return false
}

func (dc *DataController) LoadCountry(c string) bool {
	if !dc.countryExists(c) {
		return false
	}

	confirmed, err1 := covidService.CountryStatusConfirmed(c)

	if err1 != nil {
		dc.confirmedGraph.SetErrorMessage(err1.Error())
		dc.confirmedGraph.SetIsError(true)
		dc.confirmedNewGraph.SetIsError(true)
	} else {
		dc.confirmedGraph.SetIsError(false)
		dc.confirmedGraph.SetData(confirmed)

		confirmedNew := NewCasesFromCumulative(confirmed)
		dc.confirmedNewGraph.SetIsError(false)
		dc.confirmedNewGraph.SetData(confirmedNew)
	}

	deaths, err2 := covidService.CountryStatusDeaths(c)
	if err2 != nil {
		dc.deathsGraph.SetIsError(true)
		dc.deathsGraph.SetErrorMessage(err2.Error())
		dc.deathsNewGraph.SetIsError(true)
	} else {
		dc.deathsGraph.SetIsError(false)
		dc.deathsNewGraph.SetIsError(false)
		dc.deathsGraph.SetData(deaths)
		deathsNew := NewCasesFromCumulative(deaths)
		dc.deathsNewGraph.SetData(deathsNew)
	}

	recovered, err3 := covidService.CountryStatusRecovered(c)
	if err3 != nil {
		dc.recoveredGraph.SetIsError(true)
		dc.recoveredGraph.SetErrorMessage(err2.Error())
		dc.recoveredNewGraph.SetIsError(true)
	} else {
		dc.recoveredGraph.SetIsError(false)
		dc.recoveredNewGraph.SetIsError(false)
		dc.recoveredGraph.SetData(recovered)
		recoveredNew := NewCasesFromCumulative(recovered)
		dc.recoveredNewGraph.SetData(recoveredNew)
	}

	summary := CountrySummary(c, dc.summaries.Countries)
	dc.summaryView.countrySummary = summary

	return true
}

func (dc *DataController) Resize(w, h int) {
	graphTopY := 3
	totalH := h - graphTopY

	oneW := w / 2
	oneH := totalH / 2
	row1Y := graphTopY
	row2Y := graphTopY + totalH/2
	col1X := 0
	col2X := w / 2

	dc.summaryView.SetBounds(egg.MakeBounds(col1X, row1Y, oneW, 2*oneH))
	dc.confirmedGraph.SetBounds(egg.MakeBounds(col2X, row1Y, oneW, oneH))
	dc.confirmedNewGraph.SetBounds(egg.MakeBounds(col2X, row2Y, oneW, oneH))
	dc.deathsGraph.SetBounds(egg.MakeBounds(col1X, row1Y, oneW, oneH))
	dc.deathsNewGraph.SetBounds(egg.MakeBounds(col1X, row2Y, oneW, oneH))
	dc.recoveredGraph.SetBounds(egg.MakeBounds(col2X, row1Y, oneW, oneH))
	dc.recoveredNewGraph.SetBounds(egg.MakeBounds(col2X, row2Y, oneW, oneH))
}
