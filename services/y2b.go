package services

import (
	"log"
	"os/exec"
	"strings"

	"qmaru-api/configs"
)

var extCfg = configs.ExtCfg()
var pyext = extCfg["pyext"].(map[string]interface{})
var pypath = pyext["path"].(string)
var pyfiles = pyext["files"].(map[string]interface{})
var pyy2b = pyfiles["youtube"].(string)

// Y2BDownload 下载 Y2B 仅最佳匹配
func Y2BDownload(url string) string {
	cmd := exec.Command(pypath, pyy2b, url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("cmd.Run() failed", err)
	}
	results := string(out)
	return strings.TrimSpace(results)
}
