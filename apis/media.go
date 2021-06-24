package apis

import (
	"fmt"
	"strings"

	"qmaru-api/services"

	"github.com/gin-gonic/gin"
)

// Media Get Picture URLs
func Media(c *gin.Context) {
	mtype := c.Param("type")
	url := c.Query("url")
	switch mtype {
	case "news":
		urlMatch, urlType := services.PicURLCheck(url)
		if urlMatch {
			// 是否强制更新数据
			urlUpdate := strings.Split(url, "?")
			updateFlag := false
			if urlUpdate[len(urlUpdate)-1] == "update" {
				updateFlag = true
			}
			// 从数据库获取数据
			sources, counts := services.MediaFromDB(url)
			if len(sources) != 0 && !updateFlag {
				imgURLs := sources["source"]

				data := map[string]interface{}{
					"type":     mtype,
					"entities": imgURLs,
					"cache":    true,
				}
				DataHandler(c, 0, fmt.Sprintf("The news has a total of %d pictures", counts), data)
				// 从远程抓取数据
			} else if len(sources) == 0 || updateFlag {
				// 去掉 ?update 后缀的真实地址
				newurl := urlUpdate[0]
				imgURLs := services.PicData(urlType, newurl)
				if len(imgURLs) != 0 {
					services.Media2DB(mtype, urlType, newurl, imgURLs)
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
		sources, _ := services.MediaFromDB(url)
		if len(sources) != 0 {
			data := map[string]interface{}{
				"type":     mtype,
				"entities": sources["source"],
				"cache":    true,
			}
			DataHandler(c, 0, "This is a Twitter Video url", data)
			// 从远程抓取数据
		} else {
			tweetVideoData := services.TweetVideo(url)
			tweetVideoURL := tweetVideoData[0].(string)
			if tweetVideoURL != "" {
				services.Media2DB(mtype, "twitter.com", url, tweetVideoData)
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
		filename := services.Y2BDownload(url)
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
