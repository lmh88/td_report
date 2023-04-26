package api

import (
	"github.com/gogf/gf/frame/g"
)

type BaseController struct{}

func (c *BaseController) success(data interface{}) g.Map {
	return g.Map{
		"success":   true,
		"message":   "操作成功",
		"status":    1,
		"errorCode": 0,
		"data":      data,
	}
}

func (c *BaseController) fail(code int, msg string) g.Map {
	return g.Map{
		"success":   false,
		"message":   msg,
		"status":    0,
		"errorCode": code,
		"data":      nil,
	}
}
