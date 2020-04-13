package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/thomgray/covid/domain"
)

type CovidServce struct {
	client *CovidClient
}

func (cs *CovidServce) Init(client *CovidClient) *CovidServce {
	cs.client = client
	return cs
}

func (cs *CovidServce) CountryStatusConfirmed(country string) ([]domain.ValueByDay, error) {
	return cs.countryByStatus(country, StatusConfirmed)
}

// CountryStatusDeaths ...
func (cs *CovidServce) CountryStatusDeaths(country string) ([]domain.ValueByDay, error) {
	return cs.countryByStatus(country, StatusDeaths)
}

// CountryStatusRecovered ...
func (cs *CovidServce) CountryStatusRecovered(country string) ([]domain.ValueByDay, error) {
	return cs.countryByStatus(country, StatusRecovered)
}

func (cs *CovidServce) countryByStatus(country string, status Status) ([]domain.ValueByDay, error) {
	resp, err := cs.client.CountryStatus(country, status)
	if err != nil {
		return []domain.ValueByDay{}, err
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return []domain.ValueByDay{}, err
	}
	output = string(body)
	var data []domain.CountryDay = make([]domain.CountryDay, 0)

	json.Unmarshal(body, &data)

	var confirmedByDay []domain.ValueByDay = make([]domain.ValueByDay, 0)

	for _, r := range data {
		if r.Province == "" {
			date, _ := time.Parse(time.RFC3339, r.Date)
			confirmedByDay = append(confirmedByDay, domain.ValueByDay{Day: date, Value: r.Cases})
		}
	}
	return confirmedByDay, nil
}

// Countries ...
func (cs *CovidServce) Countries() []string {
	resp, _ := cs.client.Countries()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	output = string(body)

	var data []domain.Country = make([]domain.Country, 0)

	json.Unmarshal(body, &data)
	var countries []string = make([]string, len(data))
	for i, d := range data {
		countries[i] = d.Slug
	}
	return countries
}

// Summaries ...
func (cs *CovidServce) Summaries() (domain.SummaryResponse, error) {
	resp, err1 := cs.client.Summaries()
	if err1 != nil {
		return domain.SummaryResponse{}, err1
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return domain.SummaryResponse{}, err1
	}
	output = string(body)

	var data domain.SummaryResponse

	json.Unmarshal(body, &data)

	return data, nil
}
