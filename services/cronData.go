package services

import (
	"fmt"

	"qmaru-api/models"
)

// CronTime 获取定时任务执行的时间
func CronTime(cronType string) (cTime string) {
	sql := fmt.Sprintf("SELECT time FROM %s WHERE type=$1", models.CrondTimeTable)
	row := models.Psql.QueryOne(sql, cronType)
	row.Scan(&cTime)
	return
}
