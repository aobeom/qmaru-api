package main

import (
	"testing"

	"qmaru-api/service"
)

func TestPicdown(t *testing.T) {
	url := ""
	urls := service.PicData("", url)
	t.Log(urls)
}
