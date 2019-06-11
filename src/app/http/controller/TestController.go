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

/**
 * @api {GET} http://host/db 测试db访问
 * @apiVersion 1.0.0
 * @apiGroup db
 *
 * @apiParam {String} need_name 需求者名称-非空
 * @apiParam {String} e_mail 用户邮箱-非空邮箱格式
 * @apiParam  {String} phone 用户电话-非空
 * @apiParam {String} company_name 需求公司名称-非空
 * @apiParam  {String} needs_desc 需求描述-非空
 *
 * @apiSuccess {Object} code 返回码
 * @apiSuccess {Object} reason  中文解释
 * @apiSuccess {String[]} data  返回数据
 *
 * @apiSuccessExample {json} Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *          "code":0,
 *          "reason":"需求已经提交了，我们的工作人员会在2个工作日内和您取得联系!",
 *          "data":[]
 *      }
 */
func (self *TestController) TestDb(c *gin.Context) {
	userinfo := self.testSv.GetUsers()

	self.Response(c,conf.CodeOk,userinfo)

}

/**
 * @api {GET} http://host/redis 测试redis访问
 * @apiVersion 1.0.0
 * @apiGroup redis
 *
 * @apiParam {String} need_name 需求者名称-非空
 * @apiParam {String} e_mail 用户邮箱-非空邮箱格式
 * @apiParam  {String} phone 用户电话-非空
 * @apiParam {String} company_name 需求公司名称-非空
 * @apiParam  {String} needs_desc 需求描述-非空
 *
 * @apiSuccess {Object} code 返回码
 * @apiSuccess {Object} reason  中文解释
 * @apiSuccess {String[]} data  返回数据
 *
 * @apiSuccessExample {json} Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *          "code":0,
 *          "reason":"需求已经提交了，我们的工作人员会在2个工作日内和您取得联系!",
 *          "data":[]
 *      }
 */
func (self *TestController) TestRedis(c *gin.Context) {
	data := self.testSv.GetCacheInfo()
	self.Response(c,conf.CodeOk,data)
}
