package server

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/guuid"
	"strings"
	"td_report/app/bean"
	"td_report/app/service/report"
	"td_report/app/service/tool"
	"td_report/common/redis"
	"td_report/common/sendmsg/wechart"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// 报表检测最大的天数，超过了不再检测
var checkmaxday = g.Cfg().GetInt64("report.check_report_maxday")

type CheckReportServer struct {
	Adddataservice *report.AdddataService
	StatisService  *tool.StatisticsService
}

func NewCheckReportServer(adddataservice *report.AdddataService, statisService *tool.StatisticsService) *CheckReportServer {
	return &CheckReportServer{
		Adddataservice: adddataservice,
		StatisService:  statisService,
	}
}

func (t *CheckReportServer) CheckProfile(reportType, reportName, start, end string, profileList []string, ischeck bool) {
	var (
		err error
		da  map[string][]*bean.Result
	)

	if len(profileList) > 0 {
		da, err = t.StatisService.GetFileByProfile(reportName, reportType, start, end, profileList)
	} else {
		da, err = t.StatisService.GetFileWithProfileList(reportName, reportType, start, end)
	}

	if err != nil {
		logger.Logger.Error(err, "====get check file is error ")
		return
	}

	if len(da) > 0 {
		var mapdata = make(map[string][]*bean.Result, 0)

		if reportType == vars.SB {
			t.StatisService.ReportSbCheckRepository.ClearAllWithRange(reportName, start, end)
		} else if reportType == vars.SP {
			t.StatisService.ReportSpCheckRepository.ClearAllWithRange(reportName, start, end)
		} else if reportType == vars.SD {
			t.StatisService.ReportSdCheckRepository.ClearAllWithRange(reportName, start, end)
		} else if reportType == vars.DSP {
			t.StatisService.ReportDspCheckRepository.ClearAllWithRange(reportName, start, end)
		}
		
		for _, item := range da {
			mapdata[reportType] = append(mapdata[reportType], item...)
		}

		//for date, item := range da {
		//	mapdata[reportType] = append(mapdata[reportType], item...)
		//	t.adddata(item, reportName, reportType, date, ischeck)
		//}

		for key, it := range mapdata {
			if len(it) > 0 {
				if key == vars.DSP {
					t.StatisService.ReportDspCheckRepository.Addone(reportName, it)
				} else if key == vars.SD {
					t.StatisService.ReportSdCheckRepository.Addone(reportName, it)
				} else if key == vars.SB {
					t.StatisService.ReportSbCheckRepository.Addone(reportName, it)
				} else if key == vars.SP {
					t.StatisService.ReportSpCheckRepository.Addone(reportName, it)
				}
			}
		}
	}
}

func (t *CheckReportServer) adddata(da []*bean.Result, reportName, reportType, start string, ischeck bool) {
	tempdataList := make([]*bean.Tempdata, 0)
	profileIdList := make([]string, 0)
	startdate, err := time.Parse(vars.TimeLayout, start)
	if err != nil {
		logger.Logger.Error(err, "start time change error ")
		return
	}

	daygap := g.Cfg().GetInt64("queue.daygap")
	maxRetry := g.Cfg().GetInt("queue.max_retry")

	for _, v := range da {
		key := redis.WithBatchCheck(reportType, reportName, v.StartDate)
		num := t.Adddataservice.GetCount(key, reportType, reportName, v.StartDate, v.ProfileId, v.Extrant)
		// 如果超过了最大的次数则跳过
		if num > maxRetry {
			continue
		}
		tempL := &bean.Tempdata{
			ProfileId: v.ProfileId,
			Extrant:   v.Extrant,
		}
		if v.Exists {
			modifiedtime, err := time.Parse(vars.TIMEFORMAT, v.ModifileDate)
			if err != nil {
				logger.Logger.Error(err, "time change error ")
			}
			result1 := (time.Now().Unix() - modifiedtime.Unix()) / 86400
			result2 := (modifiedtime.Unix() - startdate.Unix()) / 86400
			if result1 >= daygap && result2 < checkmaxday {
				tempL.ErrorType = 2
				tempL.FileUpdate = v.ModifileDate
				tempdataList = append(tempdataList, tempL)
				profileIdList = append(profileIdList, v.ProfileId)
			}

		} else {
			tempL.ErrorType = 1
			tempdataList = append(tempdataList, tempL)
			profileIdList = append(profileIdList, v.ProfileId)
		}
	}

	if len(profileIdList) > 0 {
		t.Adddataservice.PullData(start, reportType, reportName, tempdataList)
		//不用添加到调度，记录到日志表，直接拉取数据
		reportNameList := make([]string, 0)
		reportNameList = append(reportNameList, reportName)
		//插队
		if reportType == vars.DSP {
			groupData := t.Adddataservice.DspAdgroupData(profileIdList)
			t.Adddataservice.NewAddDspData(start, vars.DSP, 4, groupData, vars.FastQueue)
		} else {

			groupData := t.Adddataservice.PpcAdgroupData(profileIdList)
			t.Adddataservice.NewAddPpcData(start, reportType, 4, groupData, vars.FastQueue)
		}
	}
}

func (t *CheckReportServer) SendMsg(da []map[string]*bean.Result, reportType string, reportName string, daylist []string) {
	str := fmt.Sprintf("报表拉取缺少文件\n报表类型:%s,报表名称:%s \n", reportType, reportName)
	emptyNum := 0
	sendmap := make(map[string][]string, 0)
	for _, date := range daylist {
		mylist := make([]string, 0)
		for _, item := range da {
			for _, v := range item {
				if v.StartDate == date {
					if v.Exists == false {
						emptyNum++
						mylist = append(mylist, v.ProfileId)
					}
				}
			}
		}

		if len(mylist) > 0 {
			sendmap[date] = mylist
		}
	}

	if emptyNum != 0 {
		send := wechart.NewSendMsg(g.Cfg().GetString("wechat.key"), g.Cfg().GetBool("wechat.open"))
		content := ""
		// 添加批次号是为了发送的内容太多，一次发送不了，分批发送
		batch := guuid.New().String()
		send.Send("", str)
		for k, v := range sendmap {
			if len(v) > 30 {
				header := fmt.Sprintf("日期：%s\n", k) + fmt.Sprintf("批次号:%s:\nprofile 列表:[", batch)
				temp := ""
				for key, con := range v {
					temp = temp + con + ","
					content = header + fmt.Sprintf("%s", temp)
					if (key+1)%30 == 0 {
						content = strings.Trim(content, ",") + "]\n"
						fmt.Println(content)
						send.Send("", content)
						temp = ""
						time.Sleep(5 * time.Second)
					}
				}

				content = strings.Trim(content, ",") + "]\n"
				send.Send("", content)

			} else {

				content = str + fmt.Sprintf("日期：%s\n", k) + strings.Join(v, ",")
				content = strings.Trim(content, ",")
				fmt.Println(content)
				send.Send("", content)
			}

			time.Sleep(10 * time.Second)
		}

	}
}
