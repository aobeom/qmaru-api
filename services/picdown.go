package services

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"qmaru-api/configs"
	"qmaru-api/models"
	"qmaru-api/utils"

	"github.com/antchfx/htmlquery"
)

// TaskManage 控制并发
type taskManage struct {
	WG    sync.WaitGroup
	GChan chan string
	Data  []interface{}
}

// MediaFromDB 读取 media 数据
func MediaFromDB(url string) (data map[string]interface{}, counts int) {
	data = make(map[string]interface{})

	sql := fmt.Sprintf("SELECT type, website, url, source FROM %s WHERE url=$1", models.MediaInfoTable)
	row := models.Psql.QueryOne(sql, url)

	var mtype string
	var mwebsite string
	var murl string
	var msource models.QMediaArray

	row.Scan(&mtype, &mwebsite, &murl, &msource)
	if mtype != "" {
		data = map[string]interface{}{
			"type":    mtype,
			"website": mwebsite,
			"url":     murl,
			"source":  msource,
		}
	}
	counts = len(msource)
	return
}

// Media2DB 保存 media 的数据
func Media2DB(mtype, website, url string, sources interface{}) {
	sql := fmt.Sprintf("INSERT INTO %s (created_at,updated_at,type,website,url,source) VALUES ($1,$2,$3,$4,$5,$6)", models.MediaInfoTable)

	createdat := int(time.Now().Unix())
	updatedat := int(time.Now().Unix())

	var newSources models.QMediaArray
	for _, source := range sources.([]interface{}) {
		newSources = append(newSources, source.(string))
	}

	models.Psql.Exec(sql, createdat, updatedat, mtype, website, url, newSources)
}

// PicURLCheck 检查适配类型
func PicURLCheck(url string) (b bool, t string) {
	urlType := regexp.MustCompile(`https?://(.*mdpr\.jp/.*|.*ameblo\.jp/.*/entry-.*|.*thetv.jp|.*tokyopopline\.com|.*instagram.com/.*|.*hustlepress\.co\.jp|.*lineblog\.me)`)
	urlResult := urlType.FindAllString(url, -1)
	if len(urlResult) != 0 {
		b = true
		t = strings.Split(urlResult[0], "/")[2]
	} else {
		b = false
		t = ""
	}
	return
}

// setHost 返回网站的域名
func setHost(s string) (h string) {
	hosts := map[string]string{
		"mdpr.jp":           "https://mdpr.jp",
		"thetv.jp":          "https://thetv.jp",
		"tokyopopline.com":  "https://tokyopopline.com",
		"instagram.com":     "https://www.instagram.com",
		"hustlepress.co.jp": "https://hustlepress.co.jp",
		"lineblog.me":       "https://lineblog.me/",
		"ameblo":            "https://ameblo.jp/",
	}
	if _, ok := hosts[s]; ok {
		h = hosts[s]
	} else {
		h = ""
	}
	return
}

// abemaAPI ameblog 图片接口
func abemaAPI(url string) (imgs []interface{}) {
	owner := strings.Split(url, "/")[3]
	entryInfo := strings.Split(url, "-")
	entryData := entryInfo[len(entryInfo)-1]
	entryID := strings.Split(entryData, ".")[0]

	apiURL := "https://blogimgapi.ameba.jp/read_ahead/get.jsonp"
	imgPrefix := "http://stat.ameba.jp"

	headers := utils.MiniHeaders{
		"User-Agent": utils.UserAgent,
	}

	params := utils.MiniParams{
		"ameba_id": owner,
		"entry_id": entryID,
		"old":      "true",
		"sp":       "false",
	}

	// res := utils.Minireq.GetBody(apiURL, headers, params)
	res := utils.Minireq.Get(apiURL, headers, params)
	resReduce1 := strings.ReplaceAll(string(res.RawData()), "Amb.Ameblo.image.Callback(", "")
	resReduce2 := strings.ReplaceAll(resReduce1, ");", "")

	resJSON := utils.DataSuite.RawMap2Map([]byte(resReduce2))
	amebloImgList := resJSON["imgList"].([]interface{})
	for i := 0; i < len(amebloImgList); i++ {
		amebloImgInfo := amebloImgList[i].(map[string]interface{})
		amebloPageURL := amebloImgInfo["pageUrl"].(string)
		if strings.Contains(amebloPageURL, entryID) {
			imgURI := amebloImgInfo["imgUrl"].(string)
			imgURL := imgPrefix + imgURI
			imgs = append(imgs, imgURL)
		}
	}
	return
}

// mdprAPI 调用远程接口
func mdprAPI(url string) (imgs []interface{}) {
	/*
		AWS-Lambda and AWS-API-Gateway
		Use dynamic IP to prevent excessive frequency
	*/
	if strings.Contains(url, "photo") {
		return
	}
	awsCfg := configs.ExtCfg()["api-gateway"].(string)
	awsAPI := awsCfg + "/prod/mdpr?url=" + url
	realURL := strings.ReplaceAll(awsAPI, "?update", "")
	res := utils.Minireq.Get(realURL)
	imgs = utils.DataSuite.RawArray2Array(res.RawData())
	return
}

// igAPI 抓取网页版数据
func igAPI(url string) (imgs []interface{}) {
	var extCfg = configs.ExtCfg()
	var pyext = extCfg["pyext"].(map[string]interface{})
	var pypath = pyext["path"].(string)
	var pyfiles = pyext["files"].(map[string]interface{})
	var pyy2b = pyfiles["ig"].(string)
	cmd := exec.Command(pypath, pyy2b, url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("cmd.Run() failed", err)
	}
	results := strings.TrimSpace(string(out))

	urls := utils.DataSuite.RawArray2Array([]byte(results))
	return urls
}

// imgURLAnalysis 读取 img 标签
func imgURLAnalysis(imgRule string, taskM *taskManage) {
	url := <-taskM.GChan
	headers := utils.MiniHeaders{
		"User-Agent": utils.UserAgent,
	}
	res := utils.Minireq.Get(url, headers)
	doc, _ := htmlquery.Parse(strings.NewReader(string(res.RawData())))
	nodes := htmlquery.Find(doc, imgRule)
	for _, node := range nodes {
		imgSrc := htmlquery.SelectAttr(node, "src")
		imgURLRaw := strings.Split(imgSrc, "?")
		imgURL := imgURLRaw[0]
		taskM.Data = append(taskM.Data, imgURL)
	}
	taskM.WG.Done()
}

// picRuleProcess 其他新闻网站通过正则处理
func picRuleProcess(url, urlType, aRule, imgRule string) (imgs []interface{}) {
	headers := utils.MiniHeaders{
		"User-Agent": utils.UserAgent,
	}
	res := utils.Minireq.Get(url, headers)

	doc, _ := htmlquery.Parse(strings.NewReader(string(res.RawData())))
	// 使用 a 标签的网站
	if aRule != "" {
		var imgIndexURLs []string
		nodes := htmlquery.Find(doc, aRule)
		for _, node := range nodes {
			aHref := htmlquery.SelectAttr(node, "href")
			imgIndexURL := setHost(urlType) + aHref
			imgIndexURLs = append(imgIndexURLs, imgIndexURL)
		}

		taskM := new(taskManage)
		taskM.GChan = make(chan string, 4)
		taskM.WG.Add(len(imgIndexURLs))
		for _, imgIndexURL := range imgIndexURLs {
			taskM.GChan <- imgIndexURL
			go imgURLAnalysis(imgRule, taskM)
		}
		taskM.WG.Wait()
		imgs = taskM.Data

		// 只用 img 标签的网站
	} else {
		nodes := htmlquery.Find(doc, imgRule)
		for _, node := range nodes {
			var imgSrc string
			imgSrc = htmlquery.SelectAttr(node, "src")
			if imgSrc == "" {
				imgSrc = htmlquery.SelectAttr(node, "href")
			}
			imgURLRaw := strings.Split(imgSrc, "?")
			imgURL := imgURLRaw[0]
			imgs = append(imgs, imgURL)
		}
	}
	return
}

// PicData 获取图片地址
func PicData(urlType, url string) (imgs []interface{}) {
	switch {
	case strings.Contains(urlType, "mdpr"):
		imgs = mdprAPI(url)
	case strings.Contains(urlType, "ameblo"):
		imgs = abemaAPI(url)
	case strings.Contains(urlType, "instagram"):
		imgs = igAPI(url)
	case strings.Contains(urlType, "thetv"):
		aRule := "//ul[@class='list_thumbnail']/li/a[@alt]"
		imgRule := "//figure/a/img|//figure/img"
		imgs = picRuleProcess(url, urlType, aRule, imgRule)
	case strings.Contains(urlType, "tokyopopline"):
		aRule := ""
		imgRule := "//dl[@class='gallery-item']/dt/a/img"
		imgOrigin := picRuleProcess(url, urlType, aRule, imgRule)

		for _, imgOri := range imgOrigin {
			imgSmall := imgOri.(string)
			imgLargeRaw := strings.Split(imgSmall, "-")
			imgLarge := imgLargeRaw[0] + ".jpg"
			imgs = append(imgs, imgLarge)
		}
	case strings.Contains(urlType, "hustlepress"):
		aRule := ""
		imgRule := "//div[@class='post_content entry-content']/div/a"
		imgs = picRuleProcess(url, urlType, aRule, imgRule)
	case strings.Contains(urlType, "lineblog"):
		aRule := ""
		imgRule := "//div[@class='article-body-inner']//*/img"
		imgOrigin := picRuleProcess(url, urlType, aRule, imgRule)

		staticPicture := "https://scdn.line-apps.com/n/line_add_friends/btn/ja.png"
		for _, imgOri := range imgOrigin {
			imgSmall := imgOri.(string)
			if imgSmall != staticPicture {
				imgLarge := strings.ReplaceAll(imgSmall, "/small", "")
				imgs = append(imgs, imgLarge)
			}
		}
	default:
		imgs = []interface{}{}
	}
	return
}
