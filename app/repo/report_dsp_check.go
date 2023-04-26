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

var ReportDspCheckRepositorySet = wire.NewSet(wire.Struct(new(ReportDspCheckRepository), "*"))

type ReportDspCheckRepository struct{}

func NewReportDspCheckRepository() *ReportDspCheckRepository {
	return &ReportDspCheckRepository{}
}

func (t *ReportDspCheckRepository) Addone(reportName string, data []*bean.Result) {
	var exists int
	tempdata := make([]model.ReportDspCheck, 0)
	for _, v := range data {
		if v.Exists == true {
			exists = 1
		} else {
			exists = 0
		}
		uuidObj, _ := guuid.NewUUID()
		da := model.ReportDspCheck{
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
			dao.ReportDspCheck.DB.Model(dao.ReportDspCheck.Table).Data(tempdata).Insert()
			tempdata = make([]model.ReportDspCheck, 0)
		}
	}

	if len(tempdata) > 0 {
		dao.ReportDspCheck.DB.Model(dao.ReportDspCheck.Table).Data(tempdata).Insert()
	}
}

func (t *ReportDspCheckRepository) DeleteByCondition(reportName string, reportDate string,
	profileId string) error {
	_, err := dao.ReportDspCheck.DB.Model(dao.ReportDspCheck.Table).
		Where("report_name=? and report_date=? and profile_id=?", reportName, reportDate, profileId).Delete()
	return err
}

func (t *ReportDspCheckRepository) ClearAll(reportName string, reportDate string) error {
	mydate, _ := time.Parse(vars.TimeLayout, reportDate)
	realdate := mydate.Format(vars.TimeFormatTpl)
	_, err := dao.ReportDspCheck.DB.Model(dao.ReportDspCheck.Table).
		Where("report_name=? and report_date=?", reportName, realdate).Delete()
	return err
}

func (t *ReportDspCheckRepository) ClearAllWithRange(reportName string, start, end string) error {
	startdate, _ := time.Parse(vars.TimeLayout, start)
	startdateT := startdate.Format(vars.TimeFormatTpl)
	enddate, _ := time.Parse(vars.TimeLayout, end)
	enddateT := enddate.Format(vars.TimeFormatTpl)
	_, err := dao.ReportDspCheck.DB.Model(dao.ReportDspCheck.Table).
		Where("report_name=? and report_date >=? and report_date<=?", reportName, startdateT, enddateT).Delete()
	return err
}

func (t *ReportDspCheckRepository) GetDataBycondition(reportDate string, status int) ([]*model.ReportDspCheck, error) {
	var data []*model.ReportDspCheck
	if err := dao.ReportDspCheck.DB.Model(dao.ReportDspCheck.Table).
		Where("report_date>? and status=?", reportDate, status).Scan(&data); err != nil {
		return nil, err
	}

	return data, nil
}
