package test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"strconv"
	"td_report/cmd/report_job"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/common/reportsystem"
	"testing"
	"time"
)

func TestReportProduct(t *testing.T) {
	//report_job.RootCmd.SetArgs([]string{"product", "--report_type", "dsp"})

	//report_job.RootCmd.SetArgs([]string{"product", "--report_type", "sb"})

	//report_job.RootCmd.SetArgs([]string{"product", "--report_type", "sd"})

	report_job.RootCmd.SetArgs([]string{"new_product", "--report_type", "sd"})

	report_job.Execute()

	t.Log("over")
}

func TestReportConsumer(t *testing.T) {
	//report_job.RootCmd.SetArgs([]string{"dsp_consumer"})

	//report_job.RootCmd.SetArgs([]string{"sb_consumer"})

	//report_job.RootCmd.SetArgs([]string{"sd_consumer"})

	report_job.RootCmd.SetArgs([]string{"new_dsp_consumer"})

	report_job.Execute()

	t.Log("over")
}

func TestPool(t *testing.T) {
	dataList := []string{"aa", "bb", "cc", "dd"}
	//poolNum := 4
	pool := reportsystem.NewPool(10)
	for {
		for i := 0; i < 10; i++ {
			for _, data := range dataList {
				pool.Add(1)
				go func(data string, num int) {
					defer pool.Done()
					fmt.Println(data + strconv.Itoa(num))
					//time.Sleep(time.Millisecond * 500)
				}(data, i)
			}
		}

		fmt.Println("for is over")
		time.Sleep(time.Second)
		fmt.Println("goroutine over")

	}
	pool.Wait()
}

func TestChan(t *testing.T) {
	intChan := make(chan int, 2)
	go func() {
		for i := range intChan {
			fmt.Println(i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go dealChan(nil)

	time.Sleep(3 * time.Second)

}

func dealChan(intChan chan int) {
	for i := 1; i < 10; i++ {
		if intChan != nil {
			intChan <- i
		}
	}
}

func TestRabbitmq(t *testing.T) {
	cfg := g.Cfg().GetString("rabbitmq.address")
	RabbitmqClient, _ := rabbitmq.NewRabbitmq(cfg)
	t.Log(RabbitmqClient)
	//conn, err := amqp.Dial(addr)
}
