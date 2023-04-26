package repo

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/google/wire"
	"td_report/app/dao"
	"td_report/app/model"
)

var ReportBatchDetailRepositorySet = wire.NewSet(wire.Struct(new(ReportBatchDetailRepository), "*"))

type ReportBatchDetailRepository struct{}

func NewReportBatchDetailRepository() *ReportBatchDetailRepository {
	return &ReportBatchDetailRepository{}
}

func (t *ReportBatchDetailRepository) Addone(batch string, reportName string, reportType string, startDate string, endDate string, layout string) (int64, error) {
	g := model.ReportBatchDetail{
		Batch:          batch,
		ReportNameList: reportName,
		Status:         1,
		ReportType:     reportType,
		StartDate:      gtime.NewFromStrLayout(startDate, layout),
		EndDate:        gtime.NewFromStrLayout(endDate, layout),
		CreateDate:     gtime.Now(),
		UpdateDate:     gtime.Now(),
	}

	return dao.ReportBatchDetail.DB.Model(dao.ReportBatchDetail.Table).Data(g).InsertAndGetId()

}

func (t *ReportBatchDetailRepository) ChangeStatus(status int, id int64, reason string) error {
	_, err := dao.ReportBatchDetail.DB.Model(dao.ReportBatchDetail.Table).Data(g.Map{"status": status, "update_date": gtime.Now(), "reason": reason}).Where("id", id).Update()
	return err
}
