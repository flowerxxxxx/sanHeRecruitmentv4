package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

// resp.go

// Resp 返回
type Resp struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// ResponseXls 向前端返回Excel文件
// 参数 content 为上面生成的io.ReadSeeker， fileTag 为返回前端的文件名
// xls backer
func ResponseXls(c *gin.Context, content io.ReadSeeker, fileTag string) {
	fileName := fmt.Sprintf("%s%s%s.xlsx", time.Now().Format("20060102"), `-`, fileTag)
	c.Writer.Header().Add("Access-Control-Expose-Headers", "content-disposition")
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, fileName))
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Header().Add("Content-Transfer-Encoding", "binary")
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), content)
}

// ErrorResp 错误返回值
func ErrorResp(c *gin.Context, status int, msg string, data ...interface{}) {
	resp(c, status, msg, data...)
}

// SuccessResp 正确返回值
func SuccessResp(c *gin.Context, msg string, data ...interface{}) {
	resp(c, 200, msg, data...)
}

// resp 返回
func resp(c *gin.Context, status int, msg string, data ...interface{}) {
	resp := Resp{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	// 将错误代码 大于210（系统级错误）的信息录入日志
	if status >= 210 {
		go ginErrorLog(c, resp)
	}
	c.JSON(http.StatusOK, resp)
}

// 错误日志的记录函数
func ginErrorLog(c *gin.Context, resp Resp) {
	log.Println(resp,
		"\nRequestUrl:"+c.Request.RequestURI,
		"\nclientIp:"+c.ClientIP(),
		"\nrequestMethod:"+c.Request.Method,
	)
}
