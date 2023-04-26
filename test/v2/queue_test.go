package v2

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"reflect"
	"td_report/app/bean"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/consumer_server"
	"td_report/app/service/report_v2/varible"
	"td_report/cmd/report_job_v2"
	"testing"
	"time"
)

func TestDeclareQueue(t *testing.T) {

	report_job_v2.RootCmd.SetArgs([]string{"queue_declare"})
	//report_job_v2.RootCmd.SetArgs([]string{"queue_declare", "--report_type=sd", "--client_tag=c2"})

	report_job_v2.Execute()

	t.Log("over")
}

func TestCheckQueue(t *testing.T) {

	report_job_v2.RootCmd.SetArgs([]string{"check_busy"})

	report_job_v2.Execute()

	t.Log("over")
}

func TestQueueLen(t *testing.T) {
	mq  := report_v2.NewMqServer()
	count := mq.GetQueueLen("sp_profile_fast")
	t.Log(count)
}

func TestReceive(t *testing.T) {
	mq := report_v2.NewMqServer()
	mq.ReceiveMsg("sp_profile_fast", dealMsg)
	t.Log("over")
}

func dealMsg(msg []byte) error {
	time.Sleep(time.Second)
	fmt.Println(string(msg))
	return nil
}

func TestRetry(t *testing.T) {
	id := varible.ReportIdMsg{
		ReportType: "sp",
		RetryCount: map[string]int{
			"wt10s": 0,
			"wt30s": 0,
			"wt60s": 1,
		},
	}
	err := bean.ReportErr{}
	//doing := true
	consumer_server.NewConsumerTwoServer().RetryHandle(&id, true)

	//for i := 0; i < 10; i++ {
	//	for k, n := range id.RetryCount {
	//		fmt.Println(k, ":", n)
	//	}
	//}
	t.Log("over", err)
}

func TestCfg(t *testing.T) {
	c1 := g.Cfg().GetString("consumer_quantity.sp", "5")
	t.Log(reflect.TypeOf(c1), c1)
	//reflect.TypeOf(c1)
	c2 := g.Cfg().Get("consumer_quantity.sp", 5)
	t.Log(reflect.TypeOf(c2), c2)
}

func TestRandMsg(t *testing.T) {
	mq  := report_v2.NewMqServer()
	for {
		for i := 0; i < 20; i++ {
			msg, _ := json.Marshal(reportIdMsg)
			err:= mq.SendMsg(varible.ReportDefaultExchange, varible.SpReportQueue, msg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		time.Sleep(time.Second)
	}
}


func TestWaitExit(t *testing.T) {
	mq := report_v2.NewMqServer()
	queue := varible.SpReportQueue
	forever := make(chan bool)
	go func () {
		n := 1
		for {
			qLen := mq.GetQueueLen(queue)
			fmt.Println(queue, "队列检查长度", qLen)

			if qLen == 0  {
				for i := 1; i <= 6; i++ {
					time.Sleep(time.Second * 3)
					fmt.Println("n", n, "i", i)
					if mq.GetQueueLen(queue) == 0 {
						if n == 6 {
							forever <- false
						}
						n++
					} else {
						n = 1
						break
					}
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()
	go func() {
		time.Sleep(5 * time.Second)
		mq.Shutdown()
	}()
	<-forever
}

func TestInitQueue(t *testing.T) {
	//report_v2.DeclareQueue("", "")
	r1 := consumer_server.AllowRun("sb", "c2")
	fmt.Println(r1)
	t.Log("over")
}


func TestCheckQueue2(t *testing.T) {
	mq := report_v2.NewMqServer()
	res := mq.ExistQueue("sp_fail")
	fmt.Println(res)
	t.Log("over")
}

func TestGetRetry(t *testing.T) {
	res := varible.GetRetryCount("sp")
	fmt.Println(res)
	t.Log("over")
}


