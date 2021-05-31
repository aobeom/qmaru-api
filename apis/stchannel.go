package apis

import (
	"qmaru-api/services"

	"github.com/gin-gonic/gin"
)

// STchannel ST Movie
func STchannel(c *gin.Context) {
	stutime := services.CronTime("stchannel")
	stData := services.STData()
	// 从数据库获取数据
	if len(stData) != 0 {
		data := map[string]interface{}{
			"time":     stutime,
			"entities": stData,
		}
		DataHandler(c, 0, "STchannel video listing", data)
	} else {
		DataHandler(c, 1, "No listing", []interface{}{})
	}
}
