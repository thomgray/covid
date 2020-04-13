package main

import (
	"github.com/thomgray/covid/domain"
)

func NewCasesFromCumulative(cumulative []domain.ValueByDay) []domain.ValueByDay {
	res := make([]domain.ValueByDay, len(cumulative))

	for i, c := range cumulative {
		if i == 0 {
			res[i] = c
		} else {
			now := cumulative[i]
			pre := cumulative[i-1]
			new := domain.ValueByDay{
				Day:   now.Day,
				Value: now.Value - pre.Value,
			}
			res[i] = new
		}
	}

	return res
}

func CountrySummary(country string, all []domain.CountrySummary) *domain.CountrySummary {
	for _, c := range all {
		if c.Slug == country {
			return &c
		}
	}
	return nil
}
