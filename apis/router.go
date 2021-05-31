package apis

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"qmaru-api/models"
	"qmaru-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type myLoggerFormat struct{}

func (f *myLoggerFormat) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

// Logger 日志中间件
func Logger(debug bool) gin.HandlerFunc {
	logger := logrus.New()

	// 输出到文件
	if debug {
		logger.Out = os.Stdout
		logger.SetLevel(logrus.DebugLevel)
	} else {
		currentPath := utils.FileSuite.LocalPath(debug)
		logPath := filepath.Join(currentPath, "logs")
		logpath := utils.FileSuite.Create(logPath)
		accessPath := filepath.Join(logpath, "access.log")
		accessFile, _ := os.OpenFile(accessPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend)

		logger.Out = accessFile
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetFormatter(new(myLoggerFormat))

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqURI := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		logger.Infof("- %s %d %s %s %s",
			reqMethod,
			statusCode,
			clientIP,
			reqURI,
			latencyTime,
		)
	}
}

// Run 执行服务
func Run(debug bool) {
	listenAddr := "localhost:8373"

	models.DBPing()
	log.Println("Listen: " + listenAddr)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 跨域
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTION"}
	router.Use(cors.New(config))

	router.Use(gin.Recovery())
	router.Use(Logger(debug))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/media/:type", Media)
		v1.GET("/drama/:type", Drama)
		v1.GET("/stchannel", STchannel)
		v1.POST("/radiko", Radiko)
	}

	_ = router.Run(listenAddr)
}
