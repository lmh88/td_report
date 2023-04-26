package server

import (
	"github.com/gogf/gf/frame/g"
	"sync"
	"td_report/app/service/check_report_file"
	"td_report/app/service/report"
	"td_report/boot"
	"td_report/common/reporttool"
	datetool "td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

type ProductServer struct {
	AdddataService *report.AdddataService
}

func NewProductServer(adddataService *report.AdddataService) *ProductServer {
	return &ProductServer{
		AdddataService: adddataService,
	}
}

// Product 队列生产者
func (t *ProductServer) Product(productday int, reportType string, specialDay int) {
	wg := &sync.WaitGroup{}
	if reportType != "" {
		wg.Add(1)
		if reportType != vars.DSP {
			go t.product(reportType, productday, wg, specialDay)
		} else {
			go t.product(vars.DSP, productday, wg, specialDay)
		}

	} else {

		wg.Add(4)
		go t.product(vars.DSP, productday, wg, specialDay)
		reportTypeList := []string{vars.SD, vars.SP, vars.SB}
		for _, item := range reportTypeList {
			go t.product(item, productday, wg, specialDay)
		}
	}

	wg.Wait()
}

func (t *ProductServer) ProductWithDate(reportType string, startdate, enddate string) {
	wg := &sync.WaitGroup{}
	if reportType != "" {
		wg.Add(1)
		if reportType != vars.DSP {
			go t.productwithDate(reportType, startdate, enddate, wg)
		} else {
			go t.productwithDate(vars.DSP, startdate, enddate, wg)
		}

	} else {

		wg.Add(4)
		go t.productwithDate(vars.DSP, startdate, enddate, wg)
		reportTypeList := []string{vars.SD, vars.SP, vars.SB}
		for _, item := range reportTypeList {
			go t.productwithDate(item, startdate, enddate, wg)
		}
	}

	wg.Wait()
}

// GetfilerProfile 过滤黑名单profile  目前暂定2天的需要过滤，其他的不需要过滤
func (t *ProductServer) GetfilerProfile(reportType string, date string) ([]string, error) {
	key := check_report_file.GetRdsKeyOfdate(reportType, date)
	return boot.RedisCommonClient.GetClient().SMembers(key).Result()
}

func (t *ProductServer) product(reporttype string, day int, wg *sync.WaitGroup, specialDay int) {
	defer wg.Done()
	var (
		productNum       = g.Cfg().GetInt64("queue.queue_all_full")
		prefixQueue      = reporttool.GetQueuePrefix(day)
		startDay, endDay time.Time
	)

	if reporttool.CheckQueueNewbatchNum(reporttype, productNum, prefixQueue) {
		logger.Logger.Info(reporttype, " product too many ,cannot add")
		return
	}

	current := time.Now()
	startDay = current.Add(time.Duration(-1*day+1) * time.Hour * 24)
	if day == vars.ProductDay14 {
		// 去掉后7天的数据
		endDay = current.Add(time.Duration(-1*7) * time.Hour * 24)
	} else {
		endDay = current
	}
	
	if daylist, err := datetool.GetDaysWithTimeRever(startDay, endDay, vars.TimeLayout); err != nil {
		logger.Logger.Info(reporttype, "product error", err.Error())
	} else {

		profileList := make([]string, 0)
		if reporttype != vars.DSP {
			if day == specialDay {
				logger.Logger.WithFields(logger.Fields{"reportType": reporttype, "day": day}).Info("ppc get profileList filter")
				for _, reportday := range daylist {
					if profileList, err = t.GetfilerProfile(reporttype, reportday); err != nil {
						logger.Logger.WithFields(logger.Fields{"err": err.Error(), "reportType": reporttype, "day": reportday}).Error("get profilelist error ")
						profileList = make([]string, 0)
						groupData := t.AdddataService.PpcAdgroupDataFilter(profileList)
						t.AdddataService.NewAddPpcData(reportday, reporttype, 0, groupData, prefixQueue)
					} else {

						logger.Logger.WithFields(logger.Fields{"reportType": reporttype, "day": reportday}).Info("ppc get profileList filter")
						groupData := t.AdddataService.PpcAdgroupDataFilter(profileList)
						t.AdddataService.NewAddPpcData(reportday, reporttype, 0, groupData, prefixQueue)
					}
				}
			} else {

				groupData := t.AdddataService.PpcAdgroupData(profileList)
				if groupData == nil {
					logger.Logger.Info(err, "ppc get profileList groupdata error")
					return
				}
				for _, reportday := range daylist {
					t.AdddataService.NewAddPpcData(reportday, reporttype, 0, groupData, prefixQueue)
				}
			}

		} else {

			if day == specialDay {
				logger.Logger.WithFields(logger.Fields{"reportType": reporttype, "day": day}).Info("dsp get profileList filter")
				for _, reportday := range daylist {
					if profileList, err = t.GetfilerProfile(reporttype, reportday); err != nil {
						logger.Logger.WithFields(logger.Fields{"err": err.Error(), "reportType": reporttype, "day": reportday}).Error("get profilelist error ")
						profileList = make([]string, 0)
						groupData := t.AdddataService.DspAdgroupDataFilter(profileList)
						t.AdddataService.NewAddDspData(reportday, reporttype, 0, groupData, prefixQueue)

					} else {
						groupData := t.AdddataService.DspAdgroupDataFilter(profileList)
						t.AdddataService.NewAddDspData(reportday, reporttype, 0, groupData, prefixQueue)
					}
				}
			} else {

				groupData := t.AdddataService.DspAdgroupData(profileList)
				if groupData == nil {
					logger.Logger.Info(err, "dsp get profileList groupdata error")
					return
				}
				for _, reportday := range daylist {
					t.AdddataService.NewAddDspData(reportday, reporttype, 0, groupData, prefixQueue)
				}
			}
		}
	}
}

func (t *ProductServer) productwithDate(reporttype string, startdate string, enddate string, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		prefixQueue = vars.SlowQueue
	)

	if daylist, err := datetool.GetDays(startdate, enddate, vars.TimeLayout); err != nil {
		logger.Logger.Info(reporttype, "product error", err.Error())
	} else {

		profileList := make([]string, 0)
		if reporttype != vars.DSP {
			groupData := t.AdddataService.PpcAdgroupData(profileList)
			if groupData == nil {
				logger.Logger.Info(err, "ppc get profileList groupdata error")
				return
			}
			for _, reportday := range daylist {
				t.AdddataService.NewAddPpcData(reportday, reporttype, 0, groupData, prefixQueue)
			}
		} else {

			groupData := t.AdddataService.DspAdgroupData(profileList)
			if groupData == nil {
				logger.Logger.Info(err, "dsp get profileList groupdata error")
				return
			}
			for _, reportday := range daylist {
				t.AdddataService.NewAddDspData(reportday, reporttype, 0, groupData, prefixQueue)
			}
		}
	}
}
