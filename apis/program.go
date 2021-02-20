package apis

import (
	"qmaru-api/service"

	"github.com/gin-gonic/gin"
)

// Program Search Program Plan
func Program(c *gin.Context) {
	kw := c.Query("kw")
	ac := c.Query("ac")
	if kw != "" && ac != "" {
		// 从数据库获取数据
		tvdata := service.ProgramFromDB(kw, ac)
		if len(tvdata) != 0 {
			data := map[string]interface{}{
				"entities": tvdata["prog_info"],
				"cache":    true,
			}
			DataHandler(c, 0, "Program information", data)
			// 从远程抓取数据
		} else {
			tvinfo := service.YahooTV(kw, ac)
			if len(tvinfo) != 0 {
				service.Program2DB(kw, ac, tvinfo)
				data := map[string]interface{}{
					"entities": tvinfo,
				}
				DataHandler(c, 0, "Program information", data)
			} else {
				DataHandler(c, 1, "No information", []interface{}{})
			}
		}
	} else {
		DataHandler(c, 1, "Parameter error", []interface{}{})
	}
}
