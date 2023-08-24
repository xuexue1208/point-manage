package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"net/http"
)

var Ping ping

type ping struct{}

func (*ping) Ping(c *gin.Context) {
	//打印请求头
	for k, v := range c.Request.Header {
		logger.Info("请求头", k, ";value:", v)
	}
	logger.Info("User-Agent:", c.GetHeader("User-Agent"))
	//请求请求体
	body, _ := ioutil.ReadAll(c.Request.Body)
	logger.Info("请求体", string(body))
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": nil,
	})
}
