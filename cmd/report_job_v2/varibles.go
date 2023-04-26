package report_job_v2

import (
	"fmt"
	"td_report/app/repo"
	"td_report/app/service/report_v2/varible"
	"td_report/common/tool"
	"td_report/vars"
	"time"
)

var (
	StartDate      string
	EndDate        string
	ReportName     string
	ReportType     string
	ProfileId      string
	QueueType      string
	ClientTag      string
)

//检查参数client_tag，如果不正确直接panic
func checkClientTag(clientTag string) {
	if clientTag != "" {
		data := repo.NewSellerClientRepository().GetAll()
		for _, id := range data {
			if clientTag == varible.GetClientTag(id) {
				return
			}
		}
		panic(clientTag + ": is clientTag error")
	}
}

func LimitProductTime() bool {
	chinaZon, _ := tool.GetChinaZon()
	currentTime := time.Now().In(chinaZon)
	hour := currentTime.Hour()
	fmt.Println(hour)
	if productDay == vars.ProductDay2 {
		if hour > 21 || hour < 5 {
			return true
		}
	}
	return false
}
