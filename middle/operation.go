package middle

//amdin操作日志
import (
	"bytes"
	"point-manage/dao"
	"point-manage/model"
	"point-manage/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

//自己实现一个type gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

//重写Write([]byte) (int, error)
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中再写一份数据
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func Oplogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		url := c.Request.URL.Path
		username, b := c.Get("username")
		if !b {
			username = "op"
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		request := strings.Join(strings.Fields(string(bodyBytes)), "")
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		//自己实现一个type gin.ResponseWriter
		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		c.Next()
		response := writer.b.String()

		//将response 转为json
		var responseJson map[string]interface{}
		json.Unmarshal([]byte(response), &responseJson)
		//将code转为string
		codeStr := fmt.Sprintf("%v", responseJson["code"])

		//next后写入数据库
		dao.Oplogs.Create(&model.Oplogs{
			Username: username.(string),
			Method:   method,
			Url:      url,
			Ip:       c.ClientIP(),
			Request:  request,
			//Response:    response,  //所有返回的数据
			Response:    codeStr, //返回数据中的code
			CreatedTime: utils.GetNowTimeString(),
		})

	}
}
