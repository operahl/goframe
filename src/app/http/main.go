package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goframe/app/http/controller"
	middle "goframe/app/http/middleware"
	"goframe/conf"
	"goframe/dao"
	"goframe/lib/logger"
	"goframe/lib/prometheus"
	"goframe/lib/util"
	"math"
	"os"
	"runtime"
)

func init() {
	var (
		configFile = "config.toml"
	)
	if len(os.Args) >= 2 {
		configFile = os.Args[1]
	}
	if !util.FileExists(configFile) {
		panic("config file " + configFile + " no exist!")
	}
	conf.ReadCfg(configFile)

	logger.InitLog()
	dao.InitRedis()
	dao.InitDB()
}

func main() {

	NumCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(int(math.Max(float64(NumCPU-1), 1)))
	gin.SetMode(conf.Config.Server.Ginmode)
	router := gin.New()
	if conf.Config.Server.Mode != "online" {
		router.Use(gin.Recovery())
	} else {
		router.Use(middle.MyRecoveryWithWriter())
	}
	testRouter := new(controller.TestController)

	//monitor
	router.GET("/metrics", middle.NocheckToken, prometheus.Handler())

	//status
	router.GET("/status", middle.NocheckToken, testRouter.Status)

	//db
	router.GET("/db", middle.NocheckToken, testRouter.TestDb)

	//redis
	router.GET("/redis", middle.NocheckToken, testRouter.TestRedis)


	router.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{"ret": 4})
		prometheus.HttpCodeCount(c, 404)
	})
	fmt.Println("Server Port " + conf.Config.Server.Port)
	router.Run(conf.Config.Server.Port)
}
