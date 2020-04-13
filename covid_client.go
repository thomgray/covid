package main

import (
	"net/http"
	"time"
)

const (
	BaseUri string = "https://api.covid19api.com"
)

type CovidClient struct {
	client *http.Client
}

func (cc *CovidClient) Init() *CovidClient {
	cc.client = &http.Client{Timeout: 2 * time.Second}
	return cc
}

func (cc *CovidClient) CountryStatus(country string, status Status) (*http.Response, error) {
	url := BaseUri + "/total/country/" + country + "/status/" + string(status)
	resp, err := cc.client.Get(url)
	return resp, err
}

func (cc *CovidClient) Countries() (*http.Response, error) {
	url := BaseUri + "/countries"
	return cc.client.Get(url)
}

func (cc *CovidClient) Summaries() (*http.Response, error) {
	url := BaseUri + "/summary"
	return cc.client.Get(url)
}
