package apis

import (
	"fmt"
	"time"

	"qmaru-api/services"

	"github.com/gin-gonic/gin"
)

// timeCalc 计算更新时间差
func timeCalc(t string) (d int64) {
	nowUnix := time.Now().Unix()

	localTimeZone, _ := time.LoadLocation("Local")
	utimeFormat, _ := time.ParseInLocation("2006-01-02 15:04:05", t, localTimeZone)
	utimeUnix := utimeFormat.Unix()

	nextTimeUnix := utimeUnix + 14400
	d = nextTimeUnix - nowUnix
	return
}

// Drama Get JP Drama List
func Drama(c *gin.Context) {
	dtype := c.Param("type")
	allowedKey := []string{
		"time", "tvbt", "subpig", "fixsub",
	}
	allowdRes := false
	for _, key := range allowedKey {
		if dtype == key {
			allowdRes = true
		}
	}
	if allowdRes {
		// 获取执行时间和更新时间
		if dtype == "time" {
			ctime := services.CronTime("drama")
			if ctime != "" {
				countdown := timeCalc(ctime)
				data := map[string]interface{}{
					"second": countdown,
					"time":   ctime,
				}
				DataHandler(c, 0, "Drama update time", data)
			} else {
				DataHandler(c, 1, "No time data", []interface{}{})
			}
			// 获取对应组的列表
		} else {
			dramaData := services.DaramaData(dtype)
			if len(dramaData) != 0 {
				data := map[string]interface{}{
					"name":     dtype,
					"entities": dramaData,
				}
				DataHandler(c, 0, fmt.Sprintf("%s Drama List", dtype), data)
			} else {
				DataHandler(c, 1, "No drama data", []interface{}{})
			}
		}
	} else {
		DataHandler(c, 1, "The subtitle group does not exist", []interface{}{})
	}
}
