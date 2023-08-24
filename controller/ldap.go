package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Ldap ldap

type ldap struct{}

//验证ldap登录
//func (*ldap) Login(ctx *gin.Context) {
//	//定义前端传入数据结构
//	params := new(struct {
//		Username string `json:"username"`
//		Password string `json:"password"`
//	})
//	if err := ctx.ShouldBindJSON(params); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg":  err.Error(),
//			"data": nil,
//			"code": 400,
//		})
//		return
//	}
//	logger.Info(params.Username, params.Password)
//	bool, err := utils.Authenticate(params.Username, params.Password)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"msg":  err.Error(),
//			"data": nil,
//			"code": 403,
//		})
//		return
//	}
//	if !bool {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"msg":  "账号或密码错误",
//			"data": nil,
//			"code": 403,
//		})
//		return
//	}
//	token, _ := utils.JWTToken.GenToken2(params.Username)
//	b, err := dao.Role.RoleAuth(params.Username)
//	var role string
//	if b {
//		role = "data"
//	} else {
//		role = ""
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg":      "登录成功",
//		"data":     nil,
//		"code":     200,
//		"token":    token,
//		"role":     role,
//		"username": params.Username,
//	})
//}

func (*ldap) Role(ctx *gin.Context) {
	params := new(struct {
		Mobile string `json:"mobile"`
		Remark string `json:"remark"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	roleid, err := dao.Role.Create(&model.Role{
		RoleName: "data",
		Remark:   params.Remark,
		Mobile:   params.Mobile,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "权限添加成功",
		"data":   nil,
		"code":   200,
		"roleid": roleid,
	})
}

func (*ldap) Info(ctx *gin.Context) {
	role, b := ctx.Get("role")
	if b {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "获取当前用户角色",
			"data": nil,
			"code": 200,
			"role": role,
		})
		return
	}
}
