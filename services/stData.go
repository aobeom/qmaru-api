package services

import (
	"fmt"
	"log"

	"qmaru-api/models"
)

// STData 读取 STchannel 的最后 15 条数据
func STData() (data []map[string]interface{}) {
	sql := fmt.Sprintf("SELECT title,picture_url,media_url,date,path FROM %s ORDER BY date DESC LIMIT 15", models.StInfoTable)
	rows, err := models.Psql.Query(sql)
	if err != nil {
		log.Panic(err)
	}

	var stitle string
	var spurl string
	var smurl string
	var sdate string
	var spath string

	for rows.Next() {
		rows.Scan(&stitle, &spurl, &smurl, &sdate, &spath)
		data = append(data, map[string]interface{}{
			"title": stitle,
			"purl":  spurl,
			"murl":  smurl,
			"date":  sdate,
			"path":  spath,
		})
	}
	return
}
