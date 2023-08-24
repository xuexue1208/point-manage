package controller

import (
	"point-manage/config"
	"point-manage/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

var Oss oss

type oss struct{}

func (*oss) GetOssToken(ctx *gin.Context) {
	oss := config.Instance.OSSConfig
	tokens := make(map[string]interface{})

	for k, _ := range oss {
		var resp map[string]interface{}
		_ = json.Unmarshal([]byte(utils.GetPolicyToken(k)), &resp)
		tokens[k] = resp
	}
	//返回
	ctx.JSON(200, gin.H{
		"msg":  "获取token成功",
		"data": tokens,
		"code": 200,
	})
}

// callback
func (*oss) Callback(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg":  "Callback Ok",
		"code": 200,
	})
}
