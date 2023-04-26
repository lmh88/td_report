package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"td_report/app/api"
)

func init() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/ping", func(r *ghttp.Request) {
			r.Response.Write("OK！")
		})
		group.ALL("/hello", api.Hello.Index)
		group.POST("/api/report/push", api.Report.Push)
	})
}
