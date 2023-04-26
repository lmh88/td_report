package test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"io/ioutil"
	rabbitmq2 "td_report/common/rabbitmqV1"
	"td_report/common/tool"
	"td_report/vars"
	"testing"
	"time"
)

func TestToolDate(t *testing.T) {
	day := 183
	today := time.Now()
	starttime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	daylist := tool.GetBetweenWeek(starttime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	for _, item := range daylist {
		fmt.Println(item)
	}
}

func TestHealcheck(t *testing.T) {
	res, err := g.Client().SetBasicAuth("guest", "guest").SetTimeout(4 * time.Second).Get("http://127.0.0.1:15672/api/health/checks/alarms")
	if err != nil {
		t.Log(err.Error(), "-------llllll")
	} else {

		result, err := ioutil.ReadAll(res.Body)
		if err != nil {

		} else {
			fmt.Println(string(result), "gggggaaa")
		}
	}
}

func TestGetQeueuLength(t *testing.T) {
	rabbitmq, _ := rabbitmq2.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	num := rabbitmq.GetLength("report:error_info")
	fmt.Println(num)
}
