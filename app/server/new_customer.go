package server

import (
	"github.com/guonaihong/gout"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/app/service/report"
	"td_report/pkg/logger"
	"td_report/vars"
)

// CustomerServer 新用户
type CustomerServer struct {
	Customer    *report.CustomerReportService
	CheckServer *CheckReportServer
}

func NewCustomerServer() *CustomerServer {
	service := app.InitializeStatisticsService()
	profileService := app.InitializeProfileService()
	schduleRepo := repo.NewReportSchduleRepository()
	retryRepo := repo.NewReportCheckRetryDetailRepository()
	addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
	checkServer := NewCheckReportServer(addQueueService, service)
	customerService := app.InitializeCustomerReportService()
	return &CustomerServer{
		Customer:    customerService,
		CheckServer: checkServer,
	}
}

// Retry 针对status=1的情况进行重试，进入队列
func (t *CustomerServer) Retry(date string) {
	var condition = map[string]interface{}{
		"create_time >=": date,
		"status":         1,
	}

	datalist, err := t.Customer.ReportBatchRepository.GetListData(condition)
	if err != nil {
		logger.Logger.Info("new customer data not found")
		return
	}

	url := "http://127.0.0.1:8081/api/report/push"
	var code int
	var result struct {
		Success   bool        `json:"success"`
		Message   string      `json:"message"`
		Status    int         `json:"status"`
		ErrorCode int         `json:"errorCode"`
		Data      interface{} `json:"data"`
	}
	for _, item := range datalist {
		if err = gout.
			// 设置POST方法和url
			POST(url).
			//Debug(true).
			// 设置json需使用SetJSON
			SetJSON(
				item.Paramas,
			).
			BindJSON(&result).
			Code(&code).
			//结束函数
			Do(); err != nil {
			logger.Logger.Error(map[string]interface{}{
				"err":  err,
				"code": code,
			})
		}
	}
}

func (t *CustomerServer) DealNewCustomer(date string) {
	var condition = map[string]interface{}{
		"create_time >=": date,
		"is_check":       0,
	}

	datalist, err := t.Customer.ReportBatchRepository.GetList(condition)
	if err != nil {
		logger.Logger.Info("new customer data not found")
		return
	}

	if len(datalist) == 0 {
		logger.Logger.Info("no new customer data not found")
		return
	}

	var check = make(chan *bean.ReportData, 10)
	var done = make(chan struct{})

	go func() {
		defer close(done)
		for v := range check {
			reportNameList := vars.ReportList[v.ReportType]
			for _, reportname := range reportNameList {
				t.CheckServer.CheckProfile(v.ReportType, reportname, v.StartDate, v.EndDate, v.Profileids, false)
			}

			t.Customer.ReportBatchRepository.ChangeCheck(1, v.Batch)
		}
	}()

	for _, item := range datalist {
		t.Customer.Report(item)
		check <- item
	}

	close(check)
	<-done
	logger.Info("task run end ")
}
