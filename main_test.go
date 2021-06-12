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
