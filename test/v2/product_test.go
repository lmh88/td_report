package v2

import (
	"fmt"
	"td_report/app/service/report_v2/product_server"
	"td_report/app/service/report_v2/varible"
	"td_report/cmd/report_job_v2"
	"td_report/cmd/tools"
	"testing"
	"time"
)

func TestProduct(t *testing.T) {

	//report_job_v2.RootCmd.SetArgs([]string{"product", "--report_type=sp", "--day=14"})
	report_job_v2.RootCmd.SetArgs([]string{"product", "--report_type=sp", "--day=2"})

	report_job_v2.Execute()

	t.Log("over")
}

func TestProductTool(t *testing.T) {

	report_job_v2.RootCmd.SetArgs([]string{"product_tool", "--report_type=sd", "--start_date=20220427", "--end_date=20220427"})

	report_job_v2.Execute()

	t.Log("over")
}

func TestProductTool2(t *testing.T) {

	report_job_v2.RootCmd.SetArgs([]string{"product_tool", "--report_type=sb", "--start_date=20220918", "--end_date=20220918", "--profile_id=1570844297534331"})

	report_job_v2.Execute()

	t.Log("over")
}


func TestGetMonth(t *testing.T) {
	//r1 := time.Now().Day()
	//time.Parse("", )

	r1, _ := time.Parse("2006-01-02", "2022-12-15")
	fmt.Println(r1)
	year, month, day := r1.Date() // time.Now().Date()

	thisM := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	nextM := thisM.AddDate(0, 1, 0)
	lastM := thisM.AddDate(0, -1, 0)

	//fmt.Println(r2, r3, r4)
	fmt.Println(day, thisM, nextM, lastM)

	r2, r3 := product_server.NewProductServer().GetLastMonth()
	fmt.Println(r2, r3)
	t.Log("over TestGetMonth")
}

func TestProductMonth(t *testing.T) {

	tools.RootCmd.SetArgs([]string{"product_month", "--report_type=sd"})

	tools.Execute()

	t.Log("over")
}


func TestConfigVar(t *testing.T) {
	fmt.Println(varible.LimitRetryQueueMap, varible.LimitReportQueueMap)
	fmt.Println(report_job_v2.LimitProductTime())
	t.Log("over")
}



