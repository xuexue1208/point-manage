package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Recovery 中间件会恢复(recovers) 任何恐慌(panics) ,如果存在恐慌，中间件将会返回http code 500
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(
		func(c *gin.Context, recovered interface{}) {
			if err, ok := recovered.(string); ok {
				c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			}
			c.AbortWithStatus(http.StatusInternalServerError)
		})
}
