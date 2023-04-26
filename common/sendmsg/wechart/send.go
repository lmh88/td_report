package wechart

import "log"

type SendMsg struct {
	Key    string
	Isopen bool
}

func NewSendMsg(key string, isopen bool) *SendMsg {
	return &SendMsg{Key: key, Isopen: isopen}
}

// Send 目前只支持文本消息
func (t *SendMsg) Send(msgType, content string) {
	if t.Isopen == false {
		return
	}
	bot := QyBot{
		Key: t.Key,
	}
	msg := Message{
		MsgType: TextStr,
		Text: Text_{
			Content: content,
		},
	}
	send, err := bot.Send(msg)
	if err != nil {
		log.Printf("err：%v\n", err)
		return
	}
	log.Printf("send：%v\n", send)
}
