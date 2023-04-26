package api

import (
	"github.com/gogf/gf/net/ghttp"
	"td_report/common/rabbitmq"
	"td_report/vars"
)

var Hello = new(helloApi)

type helloApi struct {
	BaseController
}

// Index  测试消息队列
func (s *helloApi) Index(r *ghttp.Request) {
	if err := rabbitmq.MqServer.Publish(vars.QueueMessageLog, "log"); err != nil {
		r.Response.WriteJson(s.fail(10004, err.Error()+" ：发送消息失败"))
		return
	}
	if err := rabbitmq.MqServer.Publish(vars.QueueMessageLog, "keyword"); err != nil {
		r.Response.WriteJson(s.fail(10004, err.Error()+" ：发送消息失败"))
		return
	}
	if err := rabbitmq.MqServer.Publish(vars.QueueMessageLog, "target"); err != nil {
		r.Response.WriteJson(s.fail(10004, err.Error()+" ：发送消息失败"))
		return
	}
	r.Response.Writeln("Hello World!")
}
