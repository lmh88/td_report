package common

import (
	"fmt"
	"github.com/google/wire"
	"td_report/app/repo"
	"td_report/common/file"
	"td_report/pkg/common"
	"td_report/vars"
)

var ReportTaskServiceSet = wire.NewSet(wire.Struct(new(ReportTaskService), "*"))

type ReportTaskService struct {
	ReportTaskRepos *repo.ReportTaskRepository
}

func NewReportTaskService(reportTaskRepos *repo.ReportTaskRepository) *ReportTaskService {
	return &ReportTaskService{
		ReportTaskRepos: reportTaskRepos,
	}
}

func (t *ReportTaskService) Add(localDir, TaskType, ReportType, ReportDate string) {
	t.ReportTaskRepos.Add(localDir, TaskType, ReportType, ReportDate)
}

func (t *ReportTaskService) AddMutile(reportType, date string) {
	dirpath := vars.MypathMap[reportType]
	for _, reportname := range vars.ReportList[reportType] {
		localDir := fmt.Sprintf("%s/%s/%s", dirpath, reportname, date)
		if common.GetDirFileNum(localDir) > 0 {
			t.Add(localDir, reportname, reportType, date)
		}
	}
}

func (t *ReportTaskService) GetPath(reportType string, reportName string, dateStr string) string {
	return file.GetPath(reportType, reportName, dateStr)
}
