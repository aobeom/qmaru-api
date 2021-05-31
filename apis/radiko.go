package apis

import (
	"fmt"

	"qmaru-api/services"

	"github.com/gin-gonic/gin"
)

type radioInfo struct {
	Station string `json:"station"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
}

// Radiko Download Radiko
func Radiko(c *gin.Context) {
	c.Header("Content-Typ", "application/json")

	var reqData radioInfo
	_ = c.BindJSON(&reqData)

	// 构造 JSON 数据
	station := reqData.Station
	startAt := reqData.StartAt
	endAt := reqData.EndAt

	if station != "" && startAt != "" && endAt != "" {
		// 从数据库获取数据
		fileName := fmt.Sprintf("Radiko.%s.%s.%s.raw.aac", station, startAt, endAt)
		rData := services.RadioFromDB(fileName)
		if len(rData) != 0 {
			data := map[string]interface{}{
				"entities": map[string]interface{}{
					"name": rData["name"],
					"url":  rData["url"],
					"cache":    true,
				},
			}
			DataHandler(c, 0, station, data)
			// 从远程抓取数据
		} else {
			dlurl := services.RadioGet(fileName, station, startAt, endAt)
			if dlurl != "" {
				services.Radio2DB(fileName, dlurl)
				data := map[string]interface{}{
					"entities": map[string]interface{}{
						"name": fileName,
						"url":  dlurl,
					},
				}
				DataHandler(c, 0, station, data)
			} else {
				DataHandler(c, 1, "No Radio", []interface{}{})
			}
		}
	} else {
		DataHandler(c, 1, "Parameter error", []interface{}{})
	}
}
