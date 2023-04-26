package v2

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"td_report/app/repo"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/consumer_server"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
	"td_report/cmd/report_job_v2"
	"td_report/pkg/logger"
	"testing"
)

var reportIdMsg = varible.ReportIdMsg{
ReportType:   "sp",
ProfileId:    "123",
Region:       "NA",
RefreshToken: "abc",
Timestamp:    1650611245,
ReportDate:   "20220403",
BatchKey:     "report:sp:20220421:1650527292_4ItL",
TraceId:      "sp_297150943679491_20220411_keywords_1649666384_abc",
ReportId:     "cha",
ReportName:   "keywords",
RetryCount:   varible.RetryCount,
}

func TestGetStatus(t *testing.T) {

	ctx := logger.Logger.NewTraceIDContext(context.Background(), reportIdMsg.TraceId)
	res, err := create_report.GetReportStatus(&reportIdMsg, ctx)
	fmt.Println(res, err)
	t.Log("over")
}

func TestCreateReport(t *testing.T) {
	//ctx := logger.Logger.NewTraceIDContext(context.Background(), reportIdMsg.TraceId)
	//res, err := create_report.CreateReport("NA", "", "", "", "", "", "", "", "", "", ctx)
	//fmt.Println(res, err)
	t.Log("over")
}

func TestDealRedis(t *testing.T) {
	//consumer_server.NewConsumerTwoServer().DealRedis(&reportIdMsg)
	t.Log("over")
}

func TestSpConsumer(t *testing.T) {
	//report_job_v2.RootCmd.SetArgs([]string{"sp_consumer_one", "--queue_type=slow"})

	//report_job_v2.RootCmd.SetArgs([]string{"sp_consumer_one"})

	report_job_v2.RootCmd.SetArgs([]string{"sb_consumer_one", "--client_tag=c1"})

	report_job_v2.Execute()

	t.Log("over")
}

func TestSpConsumerTwo(t *testing.T) {
	//report_job_v2.RootCmd.SetArgs([]string{"sp_consumer_one", "--queue_type=slow"})

	//report_job_v2.RootCmd.SetArgs([]string{"sd_consumer_two"})

	report_job_v2.RootCmd.SetArgs([]string{"sb_consumer_two", "--client_tag=c1"})

	report_job_v2.Execute()

	t.Log("over")
}


func TestSpConsumerRetry(t *testing.T) {

	report_job_v2.RootCmd.SetArgs([]string{"sb_consumer_try", "--client_tag=c1"})

	report_job_v2.Execute()

	t.Log("over")
}


func TestScheduleUpdate(t *testing.T) {
	repo.NewReportSchduleRepository().EndSchdule("report:sp:20220401:1650362267")
	t.Log("over")
}

func TestSelectQueue(t *testing.T) {
	mq := report_v2.NewMqServer()
	name := consumer_server.NewConsumerOneServer().SelectQueue(mq, "", "sp", "")
	t.Log(name)
}

func TestRandHalf(t *testing.T) {
	//for i := 0; i < 20; i++ {
	//	//rand.NewSource(time.Now().UnixNano())
	//	time.Sleep(time.Second)
	//	fmt.Println(rand.Intn(2))
	//}
	ot := varible.OvertimeNoticeMap["sb"]
	fmt.Println(reflect.TypeOf(ot), ot)
	t.Log("over")

}

func TestTimeout(t *testing.T) {
	r1 := strings.Contains("my name is", "is")

	r2 := varible.GetQueueBusyKey("sp", "c2", "report")

	fmt.Println(r1, r2)
	t.Log("over")
}

func TestSync(t *testing.T) {
	var rwMutex	sync.RWMutex
	//var wg sync.WaitGroup
	n := 0
	for i := 0; i < 1000; i++ {
		//wg.Add(1)
		go func() {
			//defer wg.Done()
			rwMutex.Lock()
			defer rwMutex.Unlock()
			n++
			return
		}()
	}
	//wg.Wait()
	fmt.Println(n)
	t.Log("over")
}

func TestGetFailQueue(t *testing.T) {
	res := consumer_server.NewErrServer().GetAllFailQueue()
	fmt.Println(res)
	t.Log("over")
}

func TestConsumeFail(t *testing.T) {
	report_job_v2.RootCmd.SetArgs([]string{"error_consumer"})

	report_job_v2.Execute()

	t.Log("over")
}
