package middleware

import (
	"goframe/conf"
	"goframe/lib/myredis"
	"goframe/lib/util"
	"goframe/lib/prometheus"
	"time"
	"github.com/gin-gonic/gin"
)

func CheckIp() gin.HandlerFunc {
	myfunc := func(c *gin.Context) {
		fromIp := c.ClientIP()
		ipWhiteList := conf.Config.IpWhiteList.Ips
		if !util.IpInArray(fromIp, ipWhiteList) {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
	return myfunc
}

func NocheckToken(c *gin.Context) {
	c.Set("httpReqTime", time.Now())
	c.Next()
}
func CheckToken(c *gin.Context) {
	c.Set("httpReqTime", time.Now())

	token := c.Request.Header.Get("APPINFO")
	if token == "" {
		c.AbortWithStatus(401)
		prometheus.HttpCodeCount(c, 401)
		return
	}

	tokenKey := conf.TokenPrefix + token
	tokenUid := conf.TokenUidKey

	uid, _ := myredis.String(myredis.HGet(tokenKey, tokenUid))
	if uid == "" {
		c.AbortWithStatus(401)
		prometheus.HttpCodeCount(c, 401)
		return
	}
	c.Set("uid", uid)
	c.Set("token", token)
	c.Next()
}

func CheckHeader(c *gin.Context) {
	token := c.Request.Header.Get("APPINFO")
	if token == "" {
		c.JSON(200, gin.H{"ret": 3})
		c.Abort()
		return
	}
	c.Set("token", token)
	c.Next()
}
