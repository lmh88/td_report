package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/server"
	mytool "td_report/app/service/report"
	"td_report/boot"
	"td_report/common/cryptodata"
	"td_report/common/redis"
	"td_report/common/reporttool"
	"td_report/common/sendmsg/wechart"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/guonaihong/gout"
)

func TestToken(t *testing.T) {
	day, _ := tool.GetDays("20220101", "20220101", vars.TimeLayout)
	for key, item := range day {
		fmt.Println(key, item)
	}
}

func TestSendmsg(t *testing.T) {
	key := "c0337a0d-11d1-4a59-a702-b62e91caad37"
	send := wechart.NewSendMsg(key, false)
	send.Send("", "test")
}

func TestSendNotice(t *testing.T) {
	code := 0
	key := g.Cfg().GetString("common.noticekey")
	url := "http://test.amazon.sparkxmedia.com/api/v2/sellers/update-sync"
	desc := "test1"
	id := 12345
	status := 1
	time := time.Now().Unix()
	da := map[string]interface{}{
		"id":     id,
		"status": 1,
		"t":      time,
		"desc":   desc,
	}
	fmt.Println(da)
	str := fmt.Sprintf("%s%d%d%d%s", desc, id, status, time, key)
	fmt.Println(str, "===============1")
	singn := cryptodata.Sha1(str)
	fmt.Println(singn, "==============2")
	result := new(map[string]interface{})
	err := gout.
		// 设置POST方法和url
		POST(url).
		//打开debug模式
		// 设置非结构化数据到http body里面
		// 设置json需使用SetJSON
		SetWWWForm(
			gout.H{
				"id":     id,
				"desc":   desc,
				"status": status,
				"t":      time,
				"sign":   singn,
			},
		).
		BindJSON(&result).
		Code(&code).
		//结束函数
		Do()

	if err != nil {
		g.Log().Error(err)

	}
	if code != 200 {
		re, _ := json.Marshal(result)
		fmt.Println(re)
	} else {
		fmt.Println("====ok")
	}
}

func TestGetDay(t *testing.T) {
	reportBatchrepo := repo.NewReportBatchRepository()
	date := "2022-03-03"
	var condition = map[string]interface{}{
		"create_time >=": date,
		"status":         1,
	}

	if datalist, err := reportBatchrepo.GetListData(condition); err != nil {
		fmt.Println(err)
	} else {
		if datalist != nil && len(datalist) > 0 {
			for _, item := range datalist {
				fmt.Println(item.Paramas, item.Batch)
			}
		}
	}

}

func TestGetProfile(t *testing.T) {
	profileService := app.InitializeProfileService()
	profileList, err := profileService.GetProfile(1)
	if err != nil {
		fmt.Println("=====", err)
	} else {
		fmt.Println(profileList)
	}
}

func TestGettoken(t *testing.T) {
	refreshToken := "Atzr|IwEBIA82wxkLmtfdUJwHFlyyv4btgXFcutW6CdrPH8EuNalJSsyjrdQdqvY4rlp89Q2rex41JGfIHN8pC8OOul_cOV4LYUzE6KF28XoG6vt6B1PXQALl67ArUmMRaXD6HEtWlgaBiAPmMaTuJPIyekKIieQBqnnSG8Mw3vIZrVLgPc48EjNjyV1-CsfXG2Z-QD3UvkV89sB60S_DQPtCcpXXng_Fd0Qc4UUmMRv3OnrXme56plp70P7TRWDK2j3-GHtYobUrXHarEHCAkSZL97nhmvTC7ZZwUQSOGKNpSUMVBQfTzKDVNkAw4cmpudG3AWSWFFmRX-CMO4Mh-fX1PRfmd9yAOUXSAV9DL5n2pTa5nExrLvOoprwai1HILMwrKGdrWsxEHtWZxd0u0GGBEe3tRkia6TFQxtAL0fmmPh6ZYGsAv1nYpb_URwWNYrtbjls8u_ffD52J-JxI9saj4ApZHCeSnBMvCTVWM1DBcCnmE6ibNg"
	str, err := boot.Redisclient.GetClient().Get(redis.WithAccessTokenPrefix(refreshToken)).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Println("nil ===")
		} else {
			fmt.Println(err)
			fmt.Println("==========11")
		}
	} else {
		fmt.Println(str)
	}
}

func TestG(t *testing.T) {
	//token := "Atzr|IwEBIKYC2Rz5rcqDqKhJKPm5OTS1ko2HvDpzerfa6VLrxp3IEBN2SM_Y_sQB8_USB9N7K530wYTWKKcF71l2mzeBJ4AWji7Nk99JjNjexDgzBLVT6yv2ILZcqn56bGfAtXqj41fPdoyhzfgC2fejakE3Z7Jl_Vgs1fR2N72MCQQGZSXeJN2MW4IokA33XFY3giwKM7SNRQj1RD6a-GUMInE2fbTMaJr050n4D7HkObuxYWSaFZ6oE3b-KxqoRTp8t0s5gwtZmmsg1wYLKIbKFLultCEb4WuAhZzQhAAPkmGqsZY3ExvYKlKb9GoUk0cGeQ6uLH8NrcYmX0-hxTTVza_J47ex_V8k103oPxEyZVYg2bhIG90Tq7Ieduc57MrrT7jJimADbQRHIYM-BG9LEI8narXZpDWXspIRt4yiXXhslLdRrn-UgB8Mt8CBOuwGkJlzRwHJqeNYjo_CZ5Du1ZKSXa0yJ5EzJoOOvxd2-pTW37giLw"
	//aa, code, reason, err := dbp_token.GetAccessToken(token)
	//if err != nil {
	//	t.Log(err, reason)
	//} else {
	//	t.Log(aa, code, reason)
	//}
}

func TestRedisM(t *testing.T) {
	redisclient := boot.Redisclient.GetClient()
	keylist := []string{"sd:adGroups:20220219:T00020:2897901491299779", "sd:productAds:20220219:T00020:957548701597051"}
	da, err := redisclient.HMGet("dp:sd:batch_detail:20220219", keylist...).Result()
	if err != nil {
		t.Error(err)
	} else {

		if da != nil && len(da) != 0 {
			for k, v := range da {
				fmt.Println(k, v)
			}
		} else {
			t.Log("no data")
		}
	}
}

func TestK(t *testing.T) {

}

func TestH(t *testing.T) {
	var cstSh, _ = time.LoadLocation(vars.Timezon) //上海
	fmt.Println("SH : ", time.Now().In(cstSh).Format(vars.TIMEFORMAT))
}

func TestA(t *testing.T) {
	daylist, _ := tool.GetDaysPeriod("20220201", "20220228", vars.TimeLayout, 6)
	if len(daylist) == 0 {
		t.Log("error no data")
	} else {

		for _, v := range daylist {
			fmt.Println(v.StartDate, v.EndDate)
		}
	}
}

func TestS(t *testing.T) {
	daylist, _ := tool.GetDaysDsp("20220201", "20220228", vars.TimeLayout)
	if len(daylist) == 0 {
		t.Log("error no data")
	} else {

		for _, v := range daylist {
			fmt.Println(v)
		}
	}
}

func TestYs(t *testing.T) {
	endDay := time.Now()
	day := 14
	startDay := endDay.Add(time.Duration(-1*day+1) * time.Hour * 24)
	if day == vars.ProductDay14 {
		// 去掉后7天的数据
		endDay = time.Now().Add(time.Duration(-1*7) * time.Hour * 24)
	}
	daylist, err := tool.GetDaysWithTimeRever(startDay, endDay, vars.TimeLayout)
	if err != nil {

	} else {
		for _, item := range daylist {
			fmt.Println(item)
		}
	}

}

func TestV(t *testing.T) {
	service := app.InitializeStatisticsService()
	profileService := app.InitializeProfileService()
	schduleRepo := repo.NewReportSchduleRepository()
	retryRepo := repo.NewReportCheckRetryDetailRepository()
	addQueueService := mytool.NewAdddataService(profileService, schduleRepo, retryRepo)
	checkServer := server.NewCheckReportServer(addQueueService, service)
	profileList := make([]string, 0)
	checkServer.CheckProfile("dsp", "audience", "20220201", "20220202", profileList, false)
}

func TestTime(t *testing.T) {
	//date := time.Now().Format(vars.TimeFormatTpl)
	var condition = map[string]interface{}{
		"create_time >=": "2022-03-03",
	}

	da := repo.NewReportBatchRepository()
	datalist, err := da.GetList(condition)
	if err != nil {
		t.Log(err)
	} else {
		for _, item := range datalist {
			t.Log(item)
		}
	}
}

func TestLog(t *testing.T) {
	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s", "sp", "campaigns", "testprofile"))
	logger.Logger.ErrorWithContext(ctx, "hahahah")
}

func TestDate(t *testing.T) {
	var day int = 14
	var endDay time.Time
	current := time.Now()
	startDay := current.Add(time.Duration(-1*day+1) * time.Hour * 24)
	if day == vars.ProductDay14 {
		// 去掉后7天的数据
		endDay = current.Add(time.Duration(-1*7) * time.Hour * 24)
		fmt.Println(endDay.Format(vars.TimeFormatTpl))
	} else {
		endDay = current
	}

	if daylist, err := tool.GetDaysWithTime(startDay, endDay, vars.TimeLayout); err != nil {
		t.Log("error", err)
	} else {
		for _, date := range daylist {
			fmt.Println(date)
		}
	}
}

func TestTongji(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var prefixQueue string
	var tongji = make(map[string]int, 0)
	for i := 0; i < 10000; i++ {
		result := rand.Int31n(100)
		if result <= 45 {
			prefixQueue = vars.FastQueue
			tongji[prefixQueue] = tongji[prefixQueue] + 1
		} else if result <= 70 {
			prefixQueue = vars.MiddleQueue
			tongji[prefixQueue] = tongji[prefixQueue] + 1
		} else if result <= 88 {
			prefixQueue = vars.SlowQueue
			tongji[prefixQueue] = tongji[prefixQueue] + 1
		} else {
			prefixQueue = vars.BackQueue
			tongji[prefixQueue] = tongji[prefixQueue] + 1
		}
		fmt.Println(result, prefixQueue)
	}
	fmt.Println(tongji)
}

func TestHj(t *testing.T) {
	var currentTime time.Time
	fmt.Println(runtime.GOOS)
	if runtime.GOOS == "windows" { // windows 没有对应的时区文件，获取时间错误，默认获取本地时间
		fmt.Println("=====hahahah")
		local, err2 := time.LoadLocation("Local") //服务器设置的时区
		if err2 != nil {
			fmt.Println(err2)
		}
		currentTime = time.Now().In(local)
	} else {
		chinaZon, err := tool.GetChinaZon()
		if err != nil {
			logger.Logger.Info("系统获取中国时区错误Asia/Shanghai")
			return
		}
		currentTime = time.Now().In(chinaZon)
	}

	fmt.Println(currentTime.String())
	fmt.Println(currentTime.Format(vars.TIMEFORMAT))
	currentHouer := currentTime.Hour()
	fmt.Println(currentHouer)
}

func TestGsd(t *testing.T) {
	prefixQueue := vars.FastQueue
	reportType := vars.SP
	var batch string
	var numDetail int
	num, err := vars.Cache.Get(prefixQueue)
	if err != nil {
		logger.Logger.WithFields(logger.Fields{"report_type": reportType, "queueprefix": prefixQueue}).Error(err.Error())
		return
	}
	if num == nil {
		vars.Cache.Set(prefixQueue, 0, 60*time.Second)
		numDetail = 0
	} else {
		numDetail = num.(int)
	}

	//如果短期内请求失败的次数大于10次，则跳过不请求redis,
	if numDetail > 10 {
		logger.Logger.WithFields(logger.Fields{"report_type": reportType, "queueprefix": prefixQueue}).Info("tiaoguo")
	}

	if batch, err = reporttool.GetNewBatch(reportType, prefixQueue); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(batch)
	}
}

func TestMq(t *testing.T) {
	//address := g.Cfg().GetString("rabbitmq.address")
	//rabbitmq := rabbitmq.NewRabbitmq(address)
	data := map[string]interface{}{
		"testdata": "111",
		"count":    1,
	}
	client, err := boot.GetRabbitmqClient()
	defer func() {
		client.Close()
	}()
	if err != nil {
		return
	}

	go client.Receive("test", func(data []byte) error {
		fmt.Println(string(data))
		return nil
	})

	client.Send("test", data)
}

func TestFt(t *testing.T) {
	endDay := time.Now()
	fmt.Println(endDay)
	startDay := endDay.Add(time.Duration(-1*14+1) * time.Hour * 24)
	endDay = time.Now().Add(time.Duration(-1*7) * time.Hour * 24)
	list := GetDaysWithTimeRever(startDay, endDay, "20060102")
	fmt.Println(startDay)
	fmt.Println(endDay)
	fmt.Println(list)
}

func GetDaysWithTimeRever(start time.Time, end time.Time, layout string) []string {
	hours := end.Sub(start).Hours()
	days := hours / 24
	cha := int(math.Ceil(days))
	daylist := make([]string, 0)
	for i := 0; i <= cha; i++ {
		day := end.Add(24 * time.Hour * time.Duration(-1*i))
		daylist = append(daylist, day.Format(layout))
	}

	return daylist
}
