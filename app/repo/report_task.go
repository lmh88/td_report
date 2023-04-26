package repo

import (
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

const (
	InitStatus      = "initializing"
	SucceededStatus = "succeeded"
	FailedStatus    = "task_failed"
	RetryStatus     = "result_wrong"
	DumperStatus    = "dumper_running"
)

type ReportTaskRepository struct {
}

func NewReportTaskRepository() *ReportTaskRepository {
	return &ReportTaskRepository{}
}

func (t *ReportTaskRepository) GetOne(d *bean.EventDetail) (*model.ReportTask, error, int) {
	var (
		result     gdb.Record
		reportTask *model.ReportTask
		err        error
	)

	result, err = dao.ReportTask.DB.Model(dao.ReportTask.Table).
		Where("task_type = ? and report_type = ? and ad_platform = ? and report_date = ?", d.TaskType, d.ReportType, d.AdPlatform, d.ReportDate).One()
	if err != nil || result == nil {
		return nil, err, 0
	}

	if result.IsEmpty() {
		return nil, nil, 0
	}

	err = result.Struct(&reportTask)
	if err == nil {
		return reportTask, nil, 1
	} else {
		return nil, err, 0
	}
}

func (t *ReportTaskRepository) Add(localDir, TaskType, ReportType, ReportDate string) {

	event := &bean.EventDetail{
		DataPath:   localDir,
		TaskType:   TaskType,
		ReportType: ReportType,
		TaskStatus: InitStatus,
		AdPlatform: "amazon",
		ReportDate: ReportDate,
	}

	count := 0
	for {

		count++
		if count > 3 {
			break
		}
		report, err, num := t.GetOne(event)
		if num == 0 || report == nil {
			if _, err = t.Adddata(event); err != nil {
				logger.Logger.Error(map[string]interface{}{
					"flag": ReportType + " notice db error" + err.Error(),
				})
				break
			} else {

				logger.Logger.Info(map[string]interface{}{
					"flag": ReportType + " notice  success",
				})

				break
			}

		} else {

			if err != nil {
				logger.Logger.Error(map[string]interface{}{
					"flag":  "GetOne error",
					"err":   err.Error(),
					"event": event,
				})
				break

			} else {

				//判断状态为终态才操作
				if report.TaskStatus == SucceededStatus || report.TaskStatus == FailedStatus ||
					report.TaskStatus == RetryStatus || report.TaskStatus == InitStatus {
					//写入中间状态, DumperStatus
					err = t.UpdateStatus(report, InitStatus, report.DataPath)
					if err != nil {
						logger.Logger.Error(map[string]interface{}{
							"flag": "DumperStatus error",
							"t":    report,
							"err":  err.Error(),
						})

					}

					logger.Logger.Info(map[string]interface{}{
						"flag": "DumperStatus info",
					})

					break
				}
			}
		}

		time.Sleep(time.Second * time.Duration(5))
	}
}

func (t *ReportTaskRepository) Adddata(d *bean.EventDetail) (*model.ReportTask, error) {
	timePeriod := "day"
	if d.TaskType == vars.BrandMetricsMonthly {
		timePeriod = "monthly"
	} else if d.TaskType == vars.BrandMetricsWeekly {
		timePeriod = "weekly"
	} else {
		timePeriod = "day"
	}
	g := &model.ReportTask{
		ReportDate: d.ReportDate,
		IsFull:     1,
		AdPlatform: d.AdPlatform,
		ReportType: d.ReportType,
		TaskType:   d.TaskType,
		TaskStatus: d.TaskStatus,
		TimePeriod: timePeriod,
		DataPath:   d.DataPath,
		LastUpdate: gtime.Now(),
	}
	_, err := dao.ReportTask.DB.Model(dao.ReportTask.Table).Data(g).Insert()
	return g, err
}

func (t *ReportTaskRepository) UpdateStatus(h *model.ReportTask, taskStatus, path string) error {
	h.TaskStatus = taskStatus
	h.DataPath = path
	h.LastUpdate = gtime.Now()
	_, err := dao.ReportTask.DB.Model(dao.ReportTask.Table).Data(h).Where(g.Map{"id": h.Id}).Update()
	return err
}
