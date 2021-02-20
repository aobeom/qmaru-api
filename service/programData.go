package service

import (
	"fmt"
	"qmaru-api/utils"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type programJSON struct {
	Keyword   string              `bson:"keyword"`
	AreaCode  string              `bson:"area_code"`
	ProgInfo  []map[string]string `bson:"prog_info"`
	CreatedAt time.Time           `bson:"created_at"`
}

// ProgramFromDB 读取 Program 的数据
func ProgramFromDB(kw, ac string) (data map[string]interface{}) {
	programColl := DataBase.Collection("program_info")
	fData := bson.D{
		{Key: "keyword", Value: kw},
		{Key: "area_code", Value: ac},
	}
	programData := MFind(programColl, 0, 0, fData)
	if len(programData) != 0 {
		data = programData[0]
	} else {
		data = map[string]interface{}{}
	}
	return
}

// Program2DB 保存 Program 的数据
func Program2DB(kw, ac string, tvinfo []map[string]string) {
	programColl := DataBase.Collection("program_info")
	var pdata programJSON
	pdata.Keyword = kw
	pdata.AreaCode = ac
	pdata.ProgInfo = tvinfo
	pdata.CreatedAt = time.Now()
	MInsertOne(programColl, pdata)
}

// YahooTV 获取 YahooTV 数据
func YahooTV(kw, code string) (tvinfo []map[string]string) {
	yahooSite := "https://tv.yahoo.co.jp"
	url := yahooSite + "/api/adapter"

	headers := utils.MiniHeaders{
		"User-Agent":  utils.UserAgent,
		"target-api":  "mindsSiQuery",
		"target-path": "/TVWebService/V2/contents",
	}

	data := utils.MiniFormData{
		"query":              kw,
		"siTypeId":           "1 3",
		"majorGenreId":       "",
		"areaId":             code,
		"duration":           "",
		"element":            "",
		"broadCastStartDate": "",
		"broadCastEndDate":   "",
		"start":              "0",
		"results":            "10",
		"sort":               "+broadCastStartDate",
	}

	res := utils.Minireq.Post(url, headers, data)
	resJSON := res.RawJSON()

	tvData := resJSON.(map[string]interface{})

	tvResultSet := tvData["ResultSet"].(map[string]interface{})
	tvResults := tvResultSet["Result"].([]interface{})

	for _, result := range tvResults {
		tmp := make(map[string]string)
		tvRes := result.(map[string]interface{})
		startDateF := tvRes["broadCastStartDate"].(float64)
		endDateF := tvRes["broadCastEndDate"].(float64)

		startDate := strconv.FormatFloat(startDateF, 'f', -1, 64)
		endDate := strconv.FormatFloat(endDateF, 'f', -1, 64)

		tvStartDate := utils.TimeSuite.Unix2String("01/02 15:04", startDate)
		tvEndDate := utils.TimeSuite.Unix2String("15:04", endDate)

		tmp["station"] = tvRes["serviceName"].(string)
		tmp["type"] = tvRes["siTypeName"].(string)
		tmp["title"] = tvRes["title"].(string)
		tmp["url"] = fmt.Sprintf("%s/program/%s", yahooSite, tvRes["contentsId"].(string))
		tmp["date"] = fmt.Sprintf("%s ~ %s", tvStartDate, tvEndDate)

		tvinfo = append(tvinfo, tmp)
	}
	return
}
