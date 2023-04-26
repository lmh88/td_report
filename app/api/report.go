package api

import (
	"encoding/json"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/service/report_v2/product_server"
	"td_report/boot"
	"td_report/common/queue"
	"td_report/common/rabbitmq"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/guuid"
)

var Report = new(reportApi)

// ReportData 传输过来的报表数据

type reportApi struct {
	BaseController
}

func (s *reportApi) Push(r *ghttp.Request) {
	var req bean.ReportData
	var reportService = app.InitializeToolReportService()
	var daytime = 60.0
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(s.fail(1000, "参数解析错误"))
		return
	}

	if req.ReportDataType != 0 && req.ReportDataType != 1 {
		r.Response.WriteJson(s.fail(1001, "参数reportDataType错误"))
		return
	}

	if req.ReportDataType == 1 && (req.ReportName == nil || len(req.ReportName) == 0) {
		r.Response.WriteJson(s.fail(1002, "参数reportDataType错误"))
		return
	}

	if req.ReportType != vars.SP && req.ReportType != vars.DSP && req.ReportType != vars.SD && req.ReportType != vars.SB {
		r.Response.WriteJson(s.fail(1003, "参数reportType错误"+req.ReportType))
		return
	}

	if req.EndDate == "" || req.StartDate == "" || req.CallBackUrl == "" || req.ProcessId == 0 {
		r.Response.WriteJson(s.fail(1004, "参数endDate|startDate|callBackUrl|processId为空"))
		return
	}

	start, err := tool.ParaseWithLoc(req.StartDate, vars.TimeLayout)
	if err != nil {
		r.Response.WriteJson(s.fail(1005, "参数startDate错误"))
		return
	}

	end, err := tool.ParaseWithLoc(req.EndDate, vars.TimeLayout)
	if err != nil {
		r.Response.WriteJson(s.fail(1006, "参数endDate错误"))
		return
	}

	if end.After(time.Now()) {
		r.Response.WriteJson(s.fail(1012, "参数startDate|endDate错误,不能选择未来时间"))
		return
	}

	if end.Before(start) {
		r.Response.WriteJson(s.fail(1007, "参数startDate|endDate错误"))
		return
	}

	if end.Sub(start).Hours()/24 > daytime {
		r.Response.WriteJson(s.fail(1008, "参数开始和结束时间长度超过60天"))
		return
	}

	if time.Since(end).Hours()/24 > 60 {
		r.Response.WriteJson(s.fail(1009, "参数结束时间距离当前时间长度超过60天"))
		return
	}

	if len(req.Profileids) == 0 {
		r.Response.WriteJson(s.fail(1010, "profileId 参数为空"))
		return
	}

	uuidObj, _ := guuid.NewUUID()
	req.Batch = uuidObj.String()
	data, err := json.Marshal(req)
	if err != nil {
		r.Response.WriteJson(s.fail(1011, "json格式化错误"))
		return
	}

	paramas := string(data)
	reportService.ReportBatchRepository.Addone(req.Batch, paramas)
	// marketing stream
	if req.ReportType == vars.SP {
		if mq, err := boot.GetRabbitmqClient(); err == nil {
			mq.Send(vars.FeadNewClient, req.Profileids)
			mq.Close()
		} else {
			logger.Logger.Error(err.Error())
		}
	}

	if _, ok := ppcReportType[req.ReportType]; ok {
		product_server.PushNewProfile(&req)
		r.Response.WriteJson(s.success("发送成功"))
		return
	}

	queueName, err := queue.GetQueue(req.ReportType, false)
	if err != nil {
		r.Response.WriteJson(s.fail(1012, "获取队列名称错误"))
		return
	}
	if err = rabbitmq.MqServer.Publish(queueName, paramas); err != nil {
		r.Response.WriteJson(s.fail(1013, err.Error()+" ：发送消息失败"))
		return
	}
	r.Response.WriteJson(s.success("发送成功"))
}

var ppcReportType = map[string]bool{
	vars.SP: true,
	vars.SD: true,
	vars.SB: true,
}
