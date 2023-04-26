package test

import (
	"td_report/app/service/check_report_file"
	"td_report/boot"
	"td_report/cmd/tools"
	"td_report/vars"
	"testing"
	"time"
)

func TestToolStatistics(t *testing.T) {
	tools.RootCmd.SetArgs([]string{""})

	tools.Execute()

	t.Log("over")
}

func TestCheckReportFile(t *testing.T) {
	//tools.RootCmd.SetArgs([]string{"check_report_file", "--check_date", "20220322"})
	//tools.RootCmd.SetArgs([]string{"check_report_file", "--report_type", "dsp"})

	tools.RootCmd.SetArgs([]string{"check_report_file", "--report_type=sb", "--check_date=20220521"})

	tools.Execute()

	t.Log("over")
}

func TestRds(t *testing.T) {
	rds := boot.RedisCommonClient.GetClient()

	rds.SAdd("group1", []string{"nini", "bebe", "huhu", "yiyi"})
	rds.SAdd("group2", []string{"caca", "bebe", "jiji", "yiyi"})

	l1 := rds.SInter("group1", "group2").Val()
	t.Log(l1)
}


func TestDay(t *testing.T) {
	res := check_report_file.GetLatelyThreeDay(time.Now())
	t.Log(res)
}

func TestRdsKey(t *testing.T) {
	key := check_report_file.GetBackListKey("sp", time.Now())
	t.Log(key)
	date, _ := time.Parse(vars.TimeLayout, "20220323")
	res := check_report_file.ExistProfile("dsp", date, "42722172441292")
	t.Log(res)
}

func TestGetUnauth(t *testing.T) {
	//date, _ := time.Parse(vars.TimeLayout, "20220402")
	//list, err := repo.NewReportErrorRepository().GetProfileNum("sb", date)
	//if err != nil {
	//	t.Log(err.Error())
	//}
	//for _, item := range list {
	//	t.Log(item.ProfileId, item.Num)
	//}
	//t.Log(list)

	tools.RootCmd.SetArgs([]string{"get_unauth_profile", "--report_type=sd", "--need_date=20220811"})
	tools.Execute()

	t.Log("over")
}

func  TestClearDb(t *testing.T) {
	//date, _ := time.Parse(vars.TimeLayout, "20220330")
	//date, _ := time.Parse(vars.TIMEFORMAT, "2022-03-30 12:00:00")
	//err := repo.NewReportErrorRepository().ClearByDate(date)
	//t.Log(err)
	//res, err := time.ParseDuration("-2h")
	//t.Log(res, err)
	//t.Log(time.Now().Format(vars.TimeLayout))
	//clear_tool.ClearDb("error", 14)

	tools.RootCmd.SetArgs([]string{"clear_db", "--clear_type=report_detail"})
	tools.Execute()

	t.Log("over")
}

func  TestMonitorError(t *testing.T) {

	tools.RootCmd.SetArgs([]string{"monitor_error"})
	tools.Execute()

	t.Log("over")
}
