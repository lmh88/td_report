package product_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/guonaihong/gout"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/common/cryptodata"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

func PushNewProfile(reportData *bean.ReportData) error {

	reportData = filterRepeat(reportData)
	logger.Logger.Info("reportData=>", reportData.ReportName, reportData.Profileids)
	//fmt.Println(reportData.ReportName, reportData.Profileids)
	if len(reportData.Profileids) == 0 {
		CompleteBatch(reportData.Batch)
		return nil
	}
	//return nil
	tokenList, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenById(reportData.Profileids)
	if err != nil || (tokenList == nil || len(tokenList) == 0) {
		logger.Logger.Info("error1=>", err.Error())
		return err
	}

	daylist, err := tool.GetDays(reportData.StartDate, reportData.EndDate, vars.TimeLayout)
	if err != nil {
		logger.Logger.Info("get daylist error", err.Error())
		return err
	}

	for _, day := range daylist {
		addPpc(day, reportData, tokenList)
	}
	return nil
}

func filterRepeat(reportData *bean.ReportData) *bean.ReportData {
	reportNames := make([]string, 0)
	profiles := make([]string, 0)
	var key string
	rds := boot.RedisCommonClient.GetClient()
	for _, profile := range reportData.Profileids {
		exist := false
		if len(reportData.ReportName) > 0 {
			for _, reportName := range reportData.ReportName {
				key = varible.GetNewProfileKey(reportData, reportName, profile)
				nowTime := gtime.Timestamp()
				if rds.SetNX(key, nowTime, time.Hour * 24).Val() {
					reportNames = uniqueAppend(reportNames, reportName)
					exist = true
				}
			}
		} else {
			key = varible.GetNewProfileKey(reportData, "", profile)
			nowTime := gtime.Timestamp()
			if rds.SetNX(key, nowTime, time.Hour * 24).Val() {
				exist = true
			}
		}
		if exist {
			profiles = uniqueAppend(profiles, profile)
		}
	}
	reportData.ReportName = reportNames
	reportData.Profileids = profiles
	return reportData
}

func uniqueAppend(res []string, str string) []string {
	exist := false
	for _, item := range res {
		if item == str {
			exist = true
		}
	}
	if !exist {
		res = append(res, str)
	}
	return res
}

func addPpc(reportDate string, report *bean.ReportData, profileList []*bean.ProfileToken) {
	logger.Logger.Info("profileList=>", profileList)
	if len(profileList) > 0 {
		timestamp := gtime.Timestamp()
		mq := report_v2.NewMqServer()

		for _, item := range profileList {
			clientTag := varible.GetClientTag(item.Tag)
			profileMsg := varible.ProfileMsg{
				ReportType:   report.ReportType,
				ProfileId:    item.ProfileId,
				ProfileType:  item.ProfileType,
				Region:       item.Region,
				RefreshToken: item.RefreshToken,
				Timestamp:    timestamp,
				ReportDate:   reportDate,
				BatchKey:     report.Batch,
				ClientTag:    clientTag,
				ClientId:     item.ClientId,
				ClientSecret: item.ClientSecret,
			}
			logger.Logger.Info(profileMsg)
			msg, _ := json.Marshal(profileMsg)
			//推入慢队列
			queue := varible.AddQueuePre(varible.QueueMap[report.ReportType][varible.SlowLevel], clientTag)
			err := mq.SendMsg(varible.ReportDefaultExchange, queue, msg)
			if err != nil {
				logger.Logger.Error("profile推送失败：", err.Error(), "数据结构：", profileMsg)
				//重试一次
				err = mq.Reconnect().SendMsg(varible.ReportDefaultExchange, queue, msg)
				if err != nil {
					logger.Logger.Error("profile推送重试失败：", err.Error(), "数据结构：", profileMsg)
				}

			}
		}

	}
	return
}

func CompleteBatch(batch string) {
	repo := repo.NewReportBatchRepository()
	var reportData *bean.ReportData
	one, err := repo.GetOne(batch)
	if one != nil && err == nil {
		err = repo.ChangeStatus(2, batch)
		if err != nil {
			logger.Logger.Error("CompleteBatch=>", err.Error())
		}
		err = json.Unmarshal([]byte(one.Paramas), &reportData)
		if err != nil {
			logger.Logger.Error("CompleteBatch json =>", err.Error())
		}
		notice(reportData.ProcessId, "推送成功", 2, time.Now().Unix(), reportData.CallBackUrl)
	}
}

func notice(id int, desc string, status int, time int64, url string) (error, string) {
	code := 0
	key := g.Cfg().GetString("common.noticekey")
	singn := cryptodata.Sha1(fmt.Sprintf("%s%d%d%d%s", desc, id, status, time, key))
	result := new(map[string]interface{})
	logger.Logger.Info(gout.H{
		"id":     id,
		"desc":   desc,
		"status": status,
		"t":      time,
		"sign":   singn,
	})
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
		logger.Logger.Error(err)
		return err, ""
	}
	if code != 200 {
		re, _ := json.Marshal(result)
		logger.Logger.Error(re)
		return errors.New("notice_error"), string(re)
	} else {
		logger.Logger.Info("send success")
		return nil, ""
	}
}

