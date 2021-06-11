package services

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"qmaru-api/configs"
	"qmaru-api/models"
	"qmaru-api/utils"
)

// Common Header 公共头
const (
	XRadikoDevice     = "pc"
	XRadikoUser       = "dummy_user"
	XRadikoApp        = "pc_html5"
	XRadikoAppVersion = "0.0.1"
)

// DLServer 下载控制器
type DLServer struct {
	WG    sync.WaitGroup
	Gonum chan string
}

// RadioData 请求参数
type RadioData struct {
	stationID string
	startAt   string
	endAt     string
	ft        string
	to        string
	l         string
	rtype     string
}

// encodeKey 根据偏移长度生成 KEY
func encodeKey(authkey string, offset int64, length int64) (partialkey string) {
	reader := strings.NewReader(authkey)
	buff := make([]byte, length)
	_, err := reader.ReadAt(buff, offset)
	if err != nil {
		log.Panic(err)
	}
	partialkey = base64.StdEncoding.EncodeToString(buff)
	return
}

// radikoJSKey 提取 JS 的密钥
func radikoJSKey(client http.Client) (authkey string) {
	playerURL := "http://radiko.jp/apps/js/playerCommon.js"
	req, _ := http.NewRequest("GET", playerURL, nil)
	req.Header.Add("User-Agent", utils.UserAgent)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	regKeyRule := regexp.MustCompile(`[0-9a-z]{40}`)
	authkeyMap := regKeyRule.FindAllString(string(body), -1)
	authkey = authkeyMap[0]
	return
}

// radikoChunklist 提取播放地址
func radikoChunklist(playlist string) (url string) {
	regURLRule := regexp.MustCompile(`https://.*?\.m3u8`)
	urlList := regURLRule.FindAllString(playlist, -1)
	url = urlList[0]
	return
}

// radikoAAC 提取 AAC 文件的地址
func radikoAAC(m3u8 string) (urls []string) {
	regURLRule := regexp.MustCompile(`https://.*?\.aac`)
	urls = regURLRule.FindAllString(m3u8, -1)
	return
}

// radikoAuth1 获取 token / offset / length
func radikoAuth1() (token string, offset int64, length int64) {
	auth1URL := "https://radiko.jp/v2/api/auth1"
	headers := utils.MiniHeaders{
		"User-Agent":           utils.UserAgent,
		"x-radiko-device":      XRadikoDevice,
		"x-radiko-user":        XRadikoUser,
		"x-radiko-app":         XRadikoApp,
		"x-radiko-app-version": XRadikoAppVersion,
	}

	res := utils.Minireq.Get(auth1URL, headers)
	resHeader := res.RawRes.Header

	token = resHeader.Get("X-Radiko-AuthToken")
	Keyoffset := resHeader.Get("X-Radiko-Keyoffset")
	Keylength := resHeader.Get("X-Radiko-Keylength")
	offset, _ = strconv.ParseInt(Keyoffset, 10, 64)
	length, _ = strconv.ParseInt(Keylength, 10, 64)
	return
}

// radikoAuth2 获取地区代码
func radikoAuth2(token string, partialkey string) (areaid string) {
	auth2URL := "https://radiko.jp/v2/api/auth2"

	headers := utils.MiniHeaders{
		"User-Agent":          utils.UserAgent,
		"x-radiko-device":     XRadikoDevice,
		"x-radiko-user":       XRadikoUser,
		"x-radiko-authtoken":  token,
		"x-radiko-partialkey": partialkey,
	}

	res := utils.Minireq.Get(auth2URL, headers)
	areaSplit := strings.Split(string(res.RawData()), ",")
	areaid = areaSplit[0]
	return
}

// radikoHLS 获取 AAC 下载地址
func radikoHLS(token string, areaid string, radioData *RadioData) (aacURLs []string) {
	playlistURL := "https://radiko.jp/v2/api/ts/playlist.m3u8"

	headers := utils.MiniHeaders{
		"User-Agent":         utils.UserAgent,
		"X-Radiko-AuthToken": token,
		"X-Radiko-AreaId":    areaid,
	}

	params := utils.MiniParams{
		"station_id": radioData.stationID,
		"start_at":   radioData.startAt,
		"ft":         radioData.ft,
		"end_at":     radioData.endAt,
		"to":         radioData.to,
		"l":          radioData.l,
		"type":       radioData.rtype,
	}

	chunklistRes := utils.Minireq.Get(playlistURL, headers, params)
	chunklistURL := radikoChunklist(string(chunklistRes.RawData()))

	m3u8ReqHeaders := utils.MiniHeaders{
		"User-Agent": utils.UserAgent,
	}

	m3u8Res := utils.Minireq.Get(chunklistURL, m3u8ReqHeaders)
	aacURLs = radikoAAC(string(m3u8Res.RawData()))
	if len(aacURLs) == 0 {
		aacURLs = []string{}
	}
	return
}

// rThread 并发下载
func rThread(url string, ch chan []byte, dl *DLServer) {
	headers := utils.MiniHeaders{
		"User-Agent": utils.UserAgent,
	}
	res := utils.Minireq.Get(url, headers)
	dl.WG.Done()
	<-dl.Gonum
	ch <- res.RawData()
}

// rEngine 下载器
func rEngine(urls []string, savePath string) {
	aacFile, _ := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	total := float64(len(urls))
	var thread int
	dl := new(DLServer)
	if total < 16 {
		thread = int(total)
	} else {
		thread = 16
	}
	var num int

	ch := make([]chan []byte, 1024)
	dl.Gonum = make(chan string, thread)
	dl.WG.Add(len(urls))

	for i, url := range urls {
		dl.Gonum <- url
		ch[i] = make(chan []byte)
		go rThread(url, ch[i], dl)
	}
	for _, d := range ch {
		if num == int(total) {
			break
		}
		tmp := <-d
		offset, _ := aacFile.Seek(0, 2)
		aacFile.WriteAt(tmp, offset)
		num++
	}

	dl.WG.Wait()
	defer aacFile.Close()
}

// RadioFromDB 读取 Radio 的数据
func RadioFromDB(name string) (data map[string]interface{}) {
	data = make(map[string]interface{})

	sql := fmt.Sprintf("SELECT name,url FROM %s WHERE name=$1", models.RadikoInfoTable)

	row := models.Psql.QueryOne(sql, name)
	var rname string
	var rurl string
	row.Scan(&rname, &rurl)

	if rname != "" {
		data = map[string]interface{}{
			"name": rname,
			"url":  rurl,
		}
	}
	return
}

// Radio2DB 保存 Radio 的数据
func Radio2DB(name, url string) {
	sql := fmt.Sprintf("INSERT INTO %s (created_at,updated_at,name,url) VALUES ($1,$2,$3,$4)", models.RadikoInfoTable)

	createdat := int(time.Now().Unix())
	updatedat := int(time.Now().Unix())

	models.Psql.Exec(sql, createdat, updatedat, name, url)
}

// RadioGet 获取 Radio 的数据
func RadioGet(fileName, station, startAt, endAt string) (dlurl string) {
	// authkey := radikoJSKey(client)
	authkey := "bcd151073c03b352e1ef2fd66c32209da9ca0afa"

	radioData := new(RadioData)
	radioData.stationID = station
	radioData.startAt = startAt
	radioData.endAt = endAt
	radioData.ft = startAt
	radioData.to = endAt
	radioData.l = "15"
	radioData.rtype = "b"

	token, offset, length := radikoAuth1()
	partialkey := encodeKey(authkey, offset, length)
	area := radikoAuth2(token, partialkey)
	if area != "OUT" {
		aacURLs := radikoHLS(token, area, radioData)
		if len(aacURLs) != 0 {
			mediaInfo := configs.MediaCfg()
			mediaPath := mediaInfo["media_path"].(string)
			wwwPath := mediaInfo["www_path"].(string)

			savePath := filepath.Join(mediaPath, fileName)
			rEngine(aacURLs, savePath)
			dlurl = filepath.Join(wwwPath, fileName)
		} else {
			dlurl = ""
		}
	} else {
		dlurl = ""
	}
	return
}
