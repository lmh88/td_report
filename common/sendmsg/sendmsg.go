package sendmsg

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"td_report/common/sendmsg/wechart"
)

type MessageType string

const (
	// Wechat 企业微信
	Wechat MessageType = "1000"
	// Dingdng 钉钉消息
	Dingdng MessageType = "1001"

)
type SendMsg struct {
	//消息类型
	MsgType MessageType
}

func New(MsgType MessageType )*SendMsg{
	return &SendMsg{
		MsgType: MsgType,
	}
}

func (t *SendMsg) SendMsg(msg string) {
	if t.MsgType == Wechat {
		key := g.Cfg().GetString("wechat.key")
		env := g.Cfg().GetString("server.Env")
		send := wechart.NewSendMsg(key, g.Cfg().GetBool("wechat.open"))
		send.Send("", fmt.Sprintf("env:%s\n%s", env, msg))
	}

}
