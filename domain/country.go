package domain

import "time"

type CountryDay struct {
	Country     string
	CountryCode string
	Province    string
	Cases       int
	Status      string
	Date        string
}

type ValueByDay struct {
	Day   time.Time
	Value int
}

type Country struct {
	Slug string
}

type GlobalSummary struct {
	NewConfirmed   int
	TotalConfirmed int
	NewDeaths      int
	TotalDeaths    int
	NewRecovered   int
	TotalRecovered int
}

type CountrySummary struct {
	Country        string
	Slug           string
	NewConfirmed   int
	TotalConfirmed int
	NewDeaths      int
	TotalDeaths    int
	NewRecovered   int
	TotalRecovered int
}

type SummaryResponse struct {
	Global    GlobalSummary
	Countries []CountrySummary
	Date      string
}
