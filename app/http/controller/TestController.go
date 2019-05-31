package controller

import (
	"github.com/gin-gonic/gin"
	"goframe/conf"
	"goframe/service"
	//"strings"
)


type TestController struct {
	*BaseController
	testSv service.TestService
}

func init() {
	runMode := conf.Config.Server.Mode
	if runMode != "test" {
		return
	}
}

func (self *TestController) Status(c *gin.Context) {
	c.JSON(200, gin.H{"ret": conf.CodeOk})
	return
}
func (self *TestController) TestDb(c *gin.Context) {
	userinfo := self.testSv.GetUsers()

	self.Response(c,conf.CodeOk,userinfo)

}
func (self *TestController) TestRedis(c *gin.Context) {
	data := self.testSv.GetCacheInfo()
	self.Response(c,conf.CodeOk,data)
}
