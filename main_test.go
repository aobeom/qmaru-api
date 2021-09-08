package main

import (
	"testing"

	"qmaru-api/models"
	"qmaru-api/services"
)

func TestPicdown(t *testing.T) {
	data := services.PicData("mdpr", "https://mdpr.jp/news/1818636")
	var sources models.QMediaArray
	for _, d := range data {
		sources = append(sources, d.(string))
	}
	t.Log(sources)
}

func TestRadiko(t *testing.T) {
	authkey := "bcd151073c03b352e1ef2fd66c32209da9ca0afa"
	token, offset, length := services.RadikoAuth1()
	partialkey := services.EncodeKey(authkey, offset, length)
	area := services.RadikoAuth2(token, partialkey)

	radioData := new(services.RadioData)
	radioData.StationID = "TBS"
	radioData.StartAt = "20210908090000"
	radioData.EndAt = "20210908090500"
	radioData.Ft = "20210908090000"
	radioData.To = "20210908090500"
	radioData.L = "15"
	radioData.Rtype = "b"

	aacURLs := services.RadikoHLS(token, area, radioData)
	t.Log(aacURLs)
}
