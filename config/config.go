package config

import (
	"log"
	"path/filepath"
	"qmaru-api/utils"
)

var cfgRoot = "config"

func readCfg(name string) (d map[string]interface{}) {
	cfgPath := filepath.Join(utils.FileSuite.LocalPath(Deployment()), cfgRoot, name)
	if utils.FileSuite.CheckExist(cfgPath) {
		d = utils.DataSuite.RawMap2Map(utils.FileSuite.Read(cfgPath))
	} else {
		log.Panic(cfgPath + " Not Found")
	}
	return
}

// DBCfg 数据库连接配置
func DBCfg() (d map[string]interface{}) {
	d = readCfg("database.json")
	return
}

// MediaCfg 静态文件配置
func MediaCfg() (d map[string]interface{}) {
	d = readCfg("media.json")
	return
}

// TweetCfg 推认证
func TweetCfg() (d map[string]interface{}) {
	d = readCfg("tweet.json")
	return
}

// ExtCfg 外部调用配置
func ExtCfg() (d map[string]interface{}) {
	d = readCfg("extapi.json")
	return
}
