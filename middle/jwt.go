package middle

import (
	"point-manage/dao"
	"point-manage/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取Header中的Authorization
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "请求未携带token，无权限访问",
				"data": nil,
				"code": 403,
			})
			c.Abort()
			return
		}

		// parseToken 解析token包含的信息
		claims, err := utils.JWTToken.ParseToken2(token)
		if err != nil {
			//token延期错误
			if err.Error() == "TokenExpired" {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "授权已过期",
					"data": nil,
					"code": 403,
				})
				c.Abort()
				return
			}
			//其他解析错误
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  err.Error(),
				"data": nil,
				"code": 403,
			})
			c.Abort()
			return
		}
		var data map[string]interface{}
		err = json.Unmarshal([]byte(claims), &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  err.Error(),
				"data": nil,
				"code": 403,
			})
			c.Abort()
			return
		}
		//logger.Info("解析出的token信息为: ", data, data["mobile"])
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		userinfo, roleinfo, err, _ := dao.User.GetByMobile(data["mobile"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  err.Error() + "用户不存在,token涉嫌伪造",
				"data": nil,
				"code": 403,
			})
			c.Abort()
			return
		}
		
		c.Set("username", userinfo.Name)
		c.Set("mobile", userinfo.Mobile)
		c.Set("role", roleinfo.RoleName)
		c.Next()
	}

}
