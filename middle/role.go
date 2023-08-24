package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Roleif() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role.(string) != "data" {
			c.JSON(http.StatusForbidden, gin.H{
				"msg":  "权限不足",
				"data": nil,
				"code": 403,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
