package multiarmedbandit

import (
	"math"
)

type Banner interface {
	GetID() int
	GetImpressions() float64
	GetClicks() float64
}

func PickBanner(banners []Banner) int {
	var (
		banner           Banner
		totalImpressions float64
		rating           float64 = -1
	)
	for _, b := range banners {
		imp := b.GetImpressions()
		if imp == 0 {
			imp = 1
		}
		totalImpressions += imp
	}
	for _, b := range banners {
		if r := calculate(b.GetClicks(), b.GetImpressions(), totalImpressions); r > rating {
			banner = b
			rating = r
		}
	}
	return banner.GetID()
}

func calculate(clicks, impressions, totalImpressions float64) float64 {
	if impressions == 0 {
		impressions = 1
	}
	return clicks/impressions + math.Sqrt(2*math.Log(totalImpressions)/impressions)
}
