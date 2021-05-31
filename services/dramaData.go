package services

import (
	"fmt"
	"log"
	"time"

	"qmaru-api/models"
)

// getDateRange 截取当前月份
func getDateRange() (mint, maxt string) {
	year := time.Now().Year()
	month := time.Now().Month()

	var newMonth string
	if month < 10 {
		newMonth = fmt.Sprintf("%0*d", 2, month)
	} else {
		newMonth = fmt.Sprintf("%d", month)
	}
	mint = fmt.Sprintf("%d-%s-00", year, newMonth)
	maxt = fmt.Sprintf("%d-%s-31", year, newMonth)
	return
}

// DaramaData 读取 Drama 的数据
func DaramaData(dramaType string) (data []map[string]interface{}) {
	var dtitle string
	var dtype string
	var durl string
	var ddate string
	var ddlurls models.QDramaArray
	//
	// // fixsub 没有更新日期 直接倒叙返回前 15 条数据
	if dramaType == "fixsub" {
		sql := fmt.Sprintf("SELECT title,type,url,date,dlurls FROM %s WHERE type=$1 ORDER BY id DESC LIMIT 15", models.DramaInfoTable)
		rows, err := models.Psql.Query(sql, dramaType)
		if err != nil {
			log.Panic(err)
		}
		for rows.Next() {
			rows.Scan(&dtitle, &dtype, &durl, &ddate, &ddlurls)
			data = append(data, map[string]interface{}{
				"title":  dtitle,
				"type":   dtype,
				"url":    durl,
				"date":   ddate,
				"dlurls": ddlurls,
			})
		}
		// 其他返回当前月份的数据
	} else {
		start, end := getDateRange()
		sql := fmt.Sprintf("SELECT title,type,url,date,dlurls FROM %s WHERE type=$1 and date>$2 and date<$3 ORDER BY id DESC", models.DramaInfoTable)
		rows, err := models.Psql.Query(sql, dramaType, start, end)
		if err != nil {
			log.Panic(err)
		}
		for rows.Next() {
			rows.Scan(&dtitle, &dtype, &durl, &ddate, &ddlurls)
			data = append(data, map[string]interface{}{
				"title":  dtitle,
				"type":   dtype,
				"url":    durl,
				"date":   ddate,
				"dlurls": ddlurls,
			})
		}
	}
	return
}
