package apis

import (
	"fmt"
	"qmaru-api/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// Media Get Picture URLs
func Media(c *gin.Context) {
	mtype := c.Param("type")
	url := c.Query("url")
	switch mtype {
	case "news":
		urlMatch, urlType := service.PicURLCheck(url)
		if urlMatch {
			// 是否强制更新数据
			urlUpdate := strings.Split(url, "?")
			updateFlag := false
			if urlUpdate[len(urlUpdate)-1] == "update" {
				updateFlag = true
			}
			// 从数据库获取数据
			sources := service.Media2FromDB(url)
			if len(sources) != 0 && updateFlag == false {
				imgURLs := sources["source"]
				imgCounts := len(imgURLs.([]interface{}))

				data := map[string]interface{}{
					"type":     mtype,
					"entities": imgURLs,
					"cache":    true,
				}
				DataHandler(c, 0, fmt.Sprintf("The news has a total of %d pictures", imgCounts), data)
				// 从远程抓取数据
			} else if len(sources) == 0 || updateFlag == true {
				// 去掉 ?update 后缀的真实地址
				newurl := urlUpdate[0]
				imgURLs := service.PicData(urlType, newurl)
				if len(imgURLs) != 0 {
					service.Media2DB(mtype, urlType, newurl, imgURLs)
					imgCounts := len(imgURLs)

					data := map[string]interface{}{
						"type":     mtype,
						"entities": imgURLs,
					}
					DataHandler(c, 0, fmt.Sprintf("The news has a total of %d pictures", imgCounts), data)
				} else {
					DataHandler(c, 1, "This news has no pictures", []interface{}{})
				}
			}
		} else {
			DataHandler(c, 1, "The news site is not supported", []interface{}{})
		}
	case "twitter":
		// 从数据库获取数据
		sources := service.Media2FromDB(url)
		if len(sources) != 0 {
			data := map[string]interface{}{
				"type":     mtype,
				"entities": sources["source"],
				"cache":    true,
			}
			DataHandler(c, 0, "This is a Twitter Video url", data)
			// 从远程抓取数据
		} else {
			tweetVideoURL := service.TweetVideo(url)
			if tweetVideoURL != "" {
				service.Media2DB(mtype, "twitter.com", url, tweetVideoURL)
				data := map[string]interface{}{
					"type":     mtype,
					"entities": tweetVideoURL,
				}
				DataHandler(c, 0, "This is a Twitter Video url", data)
			} else {
				DataHandler(c, 1, "This url has no video", []interface{}{})
			}
		}
	case "y2b":
		filename := service.Y2BDownload(url)
		if strings.Contains(filename, ".mp4") {
			data := map[string]interface{}{
				"type":     mtype,
				"entities": filename,
			}
			DataHandler(c, 0, "This is a Y2B Video url", data)
		} else {
			DataHandler(c, 1, "This url has no video [<200M]", []interface{}{})
		}
	default:
		DataHandler(c, 1, "The type is not supported", []interface{}{})
	}
}
