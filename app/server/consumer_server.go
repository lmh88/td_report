package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"math/rand"
	"strings"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/service/common"
	"td_report/boot"
	"td_report/common/reportsystem"
	"td_report/common/reporttool"
	"td_report/pkg/all_steps"
	"td_report/pkg/dsp/dsp_all_steps"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

type ConsumerServer struct {
	ReportTaskService *common.ReportTaskService
}

func NewConsumerServer() *ConsumerServer {
	reportTaskService := app.InitializeReportTaskService()
	return &ConsumerServer{
		ReportTaskService: reportTaskService,
	}
}

func (t *ConsumerServer) GetSpFunc(profileDataStr string, reportName string, startDate string, pool *reportsystem.Pool, timebatch string) func() {
	var profileData *bean.ProfileToken
	json.Unmarshal([]byte(profileDataStr), &profileData)
	reportType := vars.SP
	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d_%s", reportType, reportName, profileData.ProfileId, time.Now().UnixNano(), timebatch))
	return func() {
		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
		result, _ := all_steps.SPAllSteps(profileData, reportName, startDate, pool, ctx)
		if result {
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
		} else {
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
		}
	}
}

func (t *ConsumerServer) GetSdFunc(profileDataStr string, reportName string, startDate string, pool *reportsystem.Pool, timebatch string) func() {
	var profileData *bean.ProfileToken
	json.Unmarshal([]byte(profileDataStr), &profileData)
	reportType := vars.SD
	return func() {
		for _, tactic := range []string{"T00020", "T00030"} {
			ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%s_%d_%s", reportType, reportName, tactic, profileData.ProfileId, time.Now().UnixNano(), timebatch))
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
			result, _ := all_steps.SDAllSteps(profileData, reportName, startDate, tactic, pool, ctx)
			if result {
				logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
			} else {
				logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
			}
		}

	}
}

func (t *ConsumerServer) GetSbFunc(profileDataStr string, reportName string, startDate string, endDate string, pool *reportsystem.Pool, timebatch string) func() {
	var profileData *bean.ProfileToken
	json.Unmarshal([]byte(profileDataStr), &profileData)
	reportType := vars.SB
	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d_%s", reportType, reportName, profileData.ProfileId, time.Now().UnixNano(), timebatch))
	return func() {
		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
		result, _ := all_steps.SbAllSteps(profileData, reportName, startDate, endDate, pool, ctx)
		if result {
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
		} else {
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
		}
	}
}

func (t *ConsumerServer) GetDspFunc(profileDataStr string, reportName string, startDate string, pool *reportsystem.Pool, timebatch string) func() {
	var profileData *bean.DspRegionProfile
	json.Unmarshal([]byte(profileDataStr), &profileData)
	reportType := vars.DSP
	return func() {
		ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d_%s", reportType, reportName, profileData.ProfileId, time.Now().UnixNano(), timebatch))
		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
		result, _ := dsp_all_steps.DspAllSteps(profileData.Region, profileData.ProfileId, reportName, startDate, pool, ctx)
		if result {
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
		} else {
			logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
		}
	}
}

// GetRandQueuePrefix 根据随机数取模拉取数据，如果简单执行对应队列则容易一直阻塞在某一种类型里面
func (t *ConsumerServer) GetRandQueuePrefix() string {
	// 根据随机数取模: fast 队列占据45 middle占据25 slow 占据18 back 占据12
	var prefixQueue string
	result := rand.Int31n(100)
	if result <= 45 {
		prefixQueue = vars.FastQueue
	} else if result <= 70 {
		prefixQueue = vars.MiddleQueue
	} else if result <= 88 {
		prefixQueue = vars.SlowQueue
	} else {
		prefixQueue = vars.BackQueue
	}

	return prefixQueue
}

// Consumer 合并消费者，暂时未验证 // todo
func (t *ConsumerServer) Consumer(reportType string, queueType string) {
	goruntimeDsp := g.Cfg().GetInt(fmt.Sprintf("common.goruntime_%s", reportType))
	var (
		batch string
		err   error
	)

	for {

		batch = ""
		if queueType != "" {
			batch, err = reporttool.GetNewBatch(reportType, queueType)
			if err != nil || batch == "" {
				for _, item := range vars.RedisQueueList {
					if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
						break
					}
				}
			}
		} else {
			// 按照快慢中备顺序来
			for _, item := range vars.RedisQueueList {
				if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
					break
				}
			}
		}

		//如果循环多次还是没有获取到批次号则中断执行
		if batch == "" {
			time.Sleep(18 * time.Second)
			break
		}
		daylist := strings.Split(batch, ":")
		date := daylist[3]
		timebatch := daylist[4]

		batchDetailList := boot.RedisCommonClient.GetClient().LRange(batch, 0, -1).Val()
		if len(batchDetailList) == 0 {
			time.Sleep(18 * time.Second)
			break
		}
		boot.RedisCommonClient.GetClient().Del(batch).Err()
		pool := reportsystem.NewPool(goruntimeDsp)
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType, "date": date}).Info("start consumer")
		for _, item := range batchDetailList {
			for _, reportname := range vars.ReportList[reportType] {
				pool.Add(1)
				go func(ppts, preportname, pdateStr string) {
					switch reportType {
					case vars.DSP:
						t.GetDspFunc(ppts, preportname, pdateStr, pool, timebatch)()
					case vars.SD:
						t.GetSdFunc(ppts, preportname, pdateStr, pool, timebatch)()
					case vars.SP:
						t.GetSpFunc(ppts, preportname, pdateStr, pool, timebatch)()
					case vars.SB:
						t.GetSbFunc(ppts, preportname, pdateStr, pdateStr, pool, timebatch)()
					}
				}(item, reportname, date)
			}
		}

		pool.Wait()
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType}).Info("end consumer")
	}
}

func (t *ConsumerServer) DspConsumer(queueType string) {
	reportType := vars.DSP
	goruntimeDsp := g.Cfg().GetInt(fmt.Sprintf("common.goruntime_%s", reportType))
	var (
		batch string
		err   error
	)

	if err != nil {
		return
	}
	for {

		batch = ""
		if queueType != "" {
			batch, err = reporttool.GetNewBatch(reportType, queueType)
			if err != nil || batch == "" {
				for _, item := range vars.RedisQueueList {
					if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
						break
					}
				}
			}
		} else {
			// 按照快慢中备顺序来
			for _, item := range vars.RedisQueueList {
				if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
					break
				}
			}
		}

		//如果循环多次还是没有获取到批次号则中断执行
		if batch == "" {
			time.Sleep(18 * time.Second)
			break
		}

		daylist := strings.Split(batch, ":")
		date := daylist[3]
		timebatch := daylist[4]

		batchDetailList := boot.RedisCommonClient.GetClient().LRange(batch, 0, -1).Val()
		if len(batchDetailList) == 0 {
			time.Sleep(18 * time.Second)
			break
		}
		boot.RedisCommonClient.GetClient().Del(batch).Err()
		pool := reportsystem.NewPool(goruntimeDsp)
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType, "date": date}).Info("start consumer")
		for _, item := range batchDetailList {
			for _, reportname := range vars.ReportList[reportType] {
				pool.Add(1)
				go func(ppts, preportname, pdateStr string) {
					var profileData *bean.DspRegionProfile
					json.Unmarshal([]byte(ppts), &profileData)
					ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d_%s", reportType, preportname, profileData.ProfileId, time.Now().UnixNano(), timebatch))
					logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
					senddata := bean.ConsumerDetail{
						CtxId:      logger.Logger.FromTraceIDContext(ctx),
						CreateTime: time.Now().Unix(),
						ReportType: reportType,
						ReportName: preportname,
						ProfileId:  profileData.ProfileId,
						ReportDate: pdateStr,
						Batch:      batch,
					}
					result, errstruct := dsp_all_steps.DspAllSteps(profileData.Region, profileData.ProfileId, preportname, pdateStr, pool, ctx)
					if result {
						senddata.Status = 1
						senddata.UpdateTime = time.Now().Unix()
						senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
					} else {
						senddata.Status = 2
						senddata.UpdateTime = time.Now().Unix()
						senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
						senddata.ErrDesc = errstruct.ErrorReason
						boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errstruct)
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
					}

					boot.RabitmqClient.SendWithNoClose(vars.ProfileidDetail, senddata)
				}(item, reportname, date)
			}
		}

		pool.Wait()
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType}).Info("end consumer")
	}

	boot.RabitmqClient.Close()
}

func (t *ConsumerServer) SbConsumer(queueType string) {
	reportType := vars.SB
	goruntime := g.Cfg().GetInt(fmt.Sprintf("common.goruntime_%s", reportType))

	var (
		batch string
		err   error
	)

	if err != nil {
		return
	}
	for {

		batch = ""
		if queueType != "" {
			batch, err = reporttool.GetNewBatch(reportType, queueType)
			if err != nil || batch == "" {
				for _, item := range vars.RedisQueueList {
					if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
						break
					}
				}
			}
		} else {
			// 按照快慢中备顺序来
			for _, item := range vars.RedisQueueList {
				if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
					break
				}
			}
		}

		//如果循环多次还是没有获取到批次号则中断执行
		if batch == "" {
			time.Sleep(18 * time.Second)
			break
		}

		daylist := strings.Split(batch, ":")
		date := daylist[3]
		timebatch := daylist[4]

		batchDetailList := boot.RedisCommonClient.GetClient().LRange(batch, 0, -1).Val()
		if len(batchDetailList) == 0 {
			time.Sleep(18 * time.Second)
			break
		}

		boot.RedisCommonClient.GetClient().Del(batch).Err()
		pool := reportsystem.NewPool(goruntime)
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType, "date": date}).Info("start consumer")
		for _, item := range batchDetailList {
			for _, reportname := range vars.ReportList[reportType] {
				pool.Add(1)
				go func(ppts, preportname, pdateStr string) {
					var profileData *bean.ProfileToken
					json.Unmarshal([]byte(ppts), &profileData)
					ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d_%s", reportType, preportname, profileData.ProfileId, time.Now().UnixNano(), timebatch))
					logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
					senddata := bean.ConsumerDetail{
						CtxId:      logger.Logger.FromTraceIDContext(ctx),
						CreateTime: time.Now().Unix(),
						ReportType: reportType,
						ReportName: preportname,
						ProfileId:  profileData.ProfileId,
						ReportDate: pdateStr,
						Batch:      batch,
					}
					result, errstruct := all_steps.SbAllSteps(profileData, preportname, pdateStr, pdateStr, pool, ctx)
					if result {
						senddata.Status = 1
						senddata.UpdateTime = time.Now().Unix()
						senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
					} else {
						senddata.Status = 2
						senddata.UpdateTime = time.Now().Unix()
						senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
						senddata.ErrDesc = errstruct.ErrorReason
						boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errstruct)
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
					}

					boot.RabitmqClient.SendWithNoClose(vars.ProfileidDetail, senddata)
				}(item, reportname, date)
			}
		}

		pool.Wait()
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType}).Info("end consumer")
	}

	boot.RabitmqClient.Close()
}

func (t *ConsumerServer) SdConsumer(queueType string) {
	reportType := vars.SD
	goruntime := g.Cfg().GetInt(fmt.Sprintf("common.goruntime_%s", reportType))

	var (
		err   error
		batch string
	)

	for {

		batch = ""
		if queueType != "" {
			batch, err = reporttool.GetNewBatch(reportType, queueType)
			if err != nil || batch == "" {
				for _, item := range vars.RedisQueueList {
					if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
						break
					}
				}
			}
		} else {
			// 按照快慢中备顺序来
			for _, item := range vars.RedisQueueList {
				if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
					break
				}
			}
		}

		//如果循环多次还是没有获取到批次号则中断执行
		if batch == "" {
			time.Sleep(18 * time.Second)
			break
		}

		daylist := strings.Split(batch, ":")
		date := daylist[3]
		timebatch := daylist[4]

		batchDetailList := boot.RedisCommonClient.GetClient().LRange(batch, 0, -1).Val()
		if len(batchDetailList) == 0 {
			time.Sleep(18 * time.Second)
			break
		}

		boot.RedisCommonClient.GetClient().Del(batch).Err()
		pool := reportsystem.NewPool(goruntime)
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType, "date": date}).Info("start consumer")
		for _, item := range batchDetailList {
			for _, reportname := range vars.ReportList[reportType] {
				for _, tactic := range []string{"T00020", "T00030"} {
					pool.Add(1)
					go func(ppts, preportname, pdateStr, ptactic string) {
						var profileData *bean.ProfileToken
						json.Unmarshal([]byte(ppts), &profileData)
						ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%s_%d_%s", reportType, preportname, ptactic, profileData.ProfileId, time.Now().UnixNano(), timebatch))
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
						senddata := bean.ConsumerDetail{
							CtxId:      logger.Logger.FromTraceIDContext(ctx),
							CreateTime: time.Now().Unix(),
							ReportType: reportType,
							ReportName: preportname,
							ProfileId:  profileData.ProfileId,
							ReportDate: pdateStr,
							Batch:      batch,
						}
						result, errstruct := all_steps.SDAllSteps(profileData, preportname, pdateStr, ptactic, pool, ctx)
						if result {
							senddata.Status = 1
							senddata.UpdateTime = time.Now().Unix()
							senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
							logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
						} else {
							senddata.Status = 2
							senddata.UpdateTime = time.Now().Unix()
							senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
							senddata.ErrDesc = errstruct.ErrorReason
							boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errstruct)
							logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
						}

						boot.RabitmqClient.SendWithNoClose(vars.ProfileidDetail, senddata)
					}(item, reportname, date, tactic)
				}
			}
		}

		pool.Wait()
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType}).Info("end consumer")
	}

	boot.RabitmqClient.Close()
}

func (t *ConsumerServer) SpConsumer(queueType string) {
	reportType := vars.SP
	goruntime := g.Cfg().GetInt(fmt.Sprintf("common.goruntime_%s", reportType))

	var (
		batch string
		err   error
	)

	for {

		batch = ""
		if queueType != "" {
			batch, err = reporttool.GetNewBatch(reportType, queueType)
			if err != nil || batch == "" {
				for _, item := range vars.RedisQueueList {
					if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
						break
					}
				}
			}
		} else {
			// 按照快慢中备顺序来
			for _, item := range vars.RedisQueueList {
				if batch, err = reporttool.GetNewBatch(reportType, item); err == nil && batch != "" {
					break
				}
			}
		}

		//如果循环多次还是没有获取到批次号则中断执行
		if batch == "" {
			time.Sleep(18 * time.Second)
			break
		}

		daylist := strings.Split(batch, ":")
		date := daylist[3]
		timebatch := daylist[4]

		batchDetailList := boot.RedisCommonClient.GetClient().LRange(batch, 0, -1).Val()
		if len(batchDetailList) == 0 {
			time.Sleep(18 * time.Second)
			break
		}

		boot.RedisCommonClient.GetClient().Del(batch).Err()
		pool := reportsystem.NewPool(goruntime)
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType, "date": date}).Info("start consumer")
		for _, item := range batchDetailList {
			for _, reportname := range vars.ReportList[reportType] {
				pool.Add(1)
				go func(ppts, preportname, pdateStr string) {
					var profileData *bean.ProfileToken
					json.Unmarshal([]byte(ppts), &profileData)
					ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d_%s", reportType, preportname, profileData.ProfileId, time.Now().UnixNano(), timebatch))
					logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute start", reportType))
					senddata := bean.ConsumerDetail{
						CtxId:      logger.Logger.FromTraceIDContext(ctx),
						CreateTime: time.Now().Unix(),
						ReportType: reportType,
						ReportName: preportname,
						ProfileId:  profileData.ProfileId,
						ReportDate: pdateStr,
						Batch:      batch,
					}
					result, errstruct := all_steps.SPAllSteps(profileData, preportname, pdateStr, pool, ctx)
					if result {
						senddata.Status = 1
						senddata.UpdateTime = time.Now().Unix()
						senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute success", reportType))
					} else {
						senddata.Status = 2
						senddata.UpdateTime = time.Now().Unix()
						senddata.CostTime = senddata.UpdateTime - senddata.CreateTime
						senddata.ErrDesc = errstruct.ErrorReason
						boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errstruct)
						logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s report execute faile", reportType))
					}

					boot.RabitmqClient.SendWithNoClose(vars.ProfileidDetail, senddata)
				}(item, reportname, date)
			}
		}

		pool.Wait()
		logger.Logger.WithFields(logger.Fields{"batch": batch, "report_type": reportType}).Info("end consumer")
	}

	boot.RabitmqClient.Close()

}
