package vars

import "github.com/gogf/gf/frame/g"

// 队列区分环境
var env = g.Cfg().GetString("server.Env")

// 队列名称相关的常量
var (
	QueueMessageSp        = env + ":td:queue:sp"
	QueueMessageSd        = env + ":td:queue:sd"
	QueueMessageSb        = env + ":td:queue:sb"
	QueueMessageDsp       = env + ":td:queue:dsp"
	QueueMessageLog       = env + ":td:queue:log"
	QueueMessageSpCommon  = env + ":td:queue:spcommon"
	QueueMessageSdCommon  = env + ":td:queue:sdcommon"
	QueueMessageSbCommon  = env + ":td:queue:sbcommon"
	QueueMessageDspCommon = env + ":td:queue:dspcommon"
)

const (
	QueueKeySp        = "queue:0"
	QueueKeySd        = "queue:1"
	QueueKeySb        = "queue:2"
	QueueKeyDsp       = "queue:3"
	QueueKeyLog       = "queue:4"
	QueueKeySpCommon  = "queuecommon:0"
	QueueKeySdCommon  = "queuecommon:1"
	QueueKeySbCommon  = "queuecommon:2"
	QueueKeyDspCommon = "queuecommon:3"
)

var QueueList []string = []string{
	QueueMessageSp, QueueMessageSd, QueueMessageSb, QueueMessageDsp, QueueMessageLog,
}
