package repo

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/google/wire"
	"td_report/app/dao"
	"td_report/app/model"
)

var ReportSchduleRepositorySet = wire.NewSet(wire.Struct(new(ReportSchduleRepository), "*"))

type ReportSchduleRepository struct{}

func NewReportSchduleRepository() *ReportSchduleRepository {
	return &ReportSchduleRepository{}
}

// AddSchdule 全部替换为redis， 保存14天
func (t *ReportSchduleRepository) AddSchdule(data []*model.ReportSchdule) error {
	_, err := dao.ReportSchdule.DB.Model(dao.ReportSchdule.Table).Data(data).Insert()
	return err
}

func (t *ReportSchduleRepository) GetLastOne(date, reportType string) (*model.ReportSchdule, error) {
	var (
		reportSchdule *model.ReportSchdule
		err           error
	)

	err = dao.ReportSchdule.DB.Model(dao.ReportSchdule.Table).
		Where("report_date =? and report_type=?", date, reportType).
		OrderDesc("create_date").Scan(&reportSchdule)
	if err != nil {
		return nil, err
	}

	return reportSchdule, nil
}

func (t *ReportSchduleRepository) DelSchdule(reportType, reportDate, createdate string) error {
	_, err := dao.ReportSchdule.DB.Model(dao.ReportSchdule.Table).Where("report_type=? and report_date=? and create_date>?", reportType, reportDate, createdate).Delete()
	return err
}

func (t *ReportSchduleRepository) DelCurrentDaySchdule(createdate string) error {
	_, err := dao.ReportSchdule.DB.Model(dao.ReportSchdule.Table).Where("create_date>?", createdate).Delete()
	return err
}

func (t *ReportSchduleRepository) EndSchdule(batch string) error {
	_, err := dao.ReportSchdule.DB.Model(dao.ReportSchdule.Table).Data(g.Map{"end_time": gtime.Now()}).Where("batch", batch).Update()
	return err
}

func (t *ReportSchduleRepository) StartSchdule(batch string) error {
	_, err := dao.ReportSchdule.DB.Model(dao.ReportSchdule.Table).Data(g.Map{"start_time": gtime.Now()}).Where("batch", batch).Update()
	return err
}
