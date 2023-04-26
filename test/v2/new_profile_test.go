package v2

import (
	"fmt"
	"github.com/gogf/guuid"
	"math/rand"
	"td_report/app/bean"
	"td_report/app/service/report_v2/product_server"
	"td_report/common/tool"
	"td_report/vars"
	"testing"
)

func TestGetBetweenWeek(t *testing.T) {
	start := "20220401"
	end := "20220512"
	res := tool.GetBetweenWeek(start, end, vars.TimeLayout)
	fmt.Println(res)
	t.Log("over")
}

func TestProductNew(t *testing.T) {
	var req bean.ReportData
	req.Batch = guuid.New().String()
	req.Profileids = []string{"4404871489220462", "2109350508556157"}
	req.ReportType = "sd"
	req.ReportName = []string{} //[]string{"word", "name"}
	req.StartDate = "20220401"
	req.EndDate = "20220510"
	req.ReportDataType = 0
	req.ProcessId = rand.Intn(99)
	req.CallBackUrl = "chang"
	//params, _ := json.Marshal(req)
	//repo.NewReportBatchRepository().Addone(req.Batch, string(params))

	product_server.PushNewProfile(&req)
	//key := varible.GetNewProfileKey(&req)
	//fmt.Println(key)
}

func TestNewProfile(t *testing.T) {

}
