package repo

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/guuid"
	"github.com/google/wire"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"td_report/vars"
	"time"
)

var ReportSbCheckRepositorySet = wire.NewSet(wire.Struct(new(ReportSbCheckRepository), "*"))

type ReportSbCheckRepository struct{}

func NewReportSbCheckRepository() *ReportSbCheckRepository {
	return &ReportSbCheckRepository{}
}

func (t *ReportSbCheckRepository) Addone(reportName string, data []*bean.Result) {
	var exists int
	tempdata := make([]model.ReportSbCheck, 0)
	for _, v := range data {
		if v.Exists == true {
			exists = 1
		} else {
			exists = 0
		}
		uuidObj, _ := guuid.NewUUID()
		da := model.ReportSbCheck{
			Id:             uuidObj.String(),
			ReportName:     reportName,
			ReportDate:     gtime.NewFromStrLayout(v.StartDate, vars.TimeLayout),
			ProfileId:      v.ProfileId,
			Filename:       v.FileName,
			FileChangedate: v.ModifileDate,
			Status:         exists,
			Createdate:     gtime.Now(),
			Updatedate:     gtime.Now(),
		}

		tempdata = append(tempdata, da)
		if len(tempdata) >= 400 {
			dao.ReportSbCheck.DB.Model(dao.ReportSbCheck.Table).Data(tempdata).Insert()
			tempdata = make([]model.ReportSbCheck, 0)
		}
	}

	if len(tempdata) > 0 {
		dao.ReportSbCheck.DB.Model(dao.ReportSbCheck.Table).Data(tempdata).Insert()
	}
}

func (t *ReportSbCheckRepository) DeleteByCondition(reportName string, reportDate string,
	profileId string) error {
	_, err := dao.ReportSbCheck.DB.Model(dao.ReportSbCheck.Table).
		Where("report_name=? and report_date=? and profile_id=?", reportName, reportDate, profileId).Delete()
	return err
}

func (t *ReportSbCheckRepository) ClearAll(reportName string, reportDate string) error {
	mydate, _ := time.Parse(vars.TimeLayout, reportDate)
	realdate := mydate.Format(vars.TimeFormatTpl)
	_, err := dao.ReportSbCheck.DB.Model(dao.ReportSbCheck.Table).
		Where("report_name=? and report_date=?", reportName, realdate).Delete()
	return err
}

func (t *ReportSbCheckRepository) ClearAllWithRange(reportName string, start, end string) error {
	startdate, _ := time.Parse(vars.TimeLayout, start)
	startdateT := startdate.Format(vars.TimeFormatTpl)
	enddate, _ := time.Parse(vars.TimeLayout, end)
	enddateT := enddate.Format(vars.TimeFormatTpl)
	_, err := dao.ReportSbCheck.DB.Model(dao.ReportSbCheck.Table).
		Where("report_name=? and report_date >=? and report_date<=?", reportName, startdateT, enddateT).Delete()
	return err
}

func (t *ReportSbCheckRepository) GetDataBycondition(reportDate string, status int) ([]*model.ReportSbCheck, error) {
	var data []*model.ReportSbCheck
	if err := dao.ReportSbCheck.DB.Model(dao.ReportSbCheck.Table).
		Where("report_date>? and status=?", reportDate, status).Scan(&data); err != nil {
		return nil, err
	}

	return data, nil
}
