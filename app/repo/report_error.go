package repo

import (
	"log"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/guuid"
	"github.com/google/wire"
)

var ReportErrorRepositorySet = wire.NewSet(wire.Struct(new(ReportErrorRepository), "*"))

type ReportErrorRepository struct{}

func NewReportErrorRepository() *ReportErrorRepository {
	return &ReportErrorRepository{}
}

func (t *ReportErrorRepository) AddOne(data *bean.ReportErr) {
	uuidObj, _ := guuid.NewUUID()
	da := model.ReportError{
		Id:          uuidObj.String(),
		ErrorReason: data.ErrorReason,
		ErrorType:   data.ErrorType,
		KeyParamas:  data.KeyParam,
		ReportType:  data.ReportType,
		ReportName:  data.ReportName,
		ReportDate:  data.ReportDate,
		ProfileId:   data.ProfileId,
		Extrant:     data.Extra,
		CreateDate:  gtime.Now(),
	}

	dao.ReportError.DB.Model(dao.ReportError.Table).Data(da).Insert()
}

const (
	ReportErrorTypeOne = 1 //第一次请求失败
	ReportErrorTypeTwo = 2 //第二次请求失败
	ReportErrorTypeThree = 3 //第三次请求失败
	ReportErrorTypeTimeOut = 4 //等待超时未完成
	ReportErrorTypeRetry = 5 //429尝试次数过多
)

type ProfileNum struct {
	Num       int    `orm:"num"  json:"num"`              //数量
	ProfileId string `orm:"profile_id"  json:"profileId"` //profile
}

func (t *ReportErrorRepository) GetProfileNum(reportType string, date time.Time, errorType int) (profiles []ProfileNum, err error) {
	startTime := date
	oneDay, _ := time.ParseDuration("24h")
	endTime := date.Add(oneDay)

	profiles = make([]ProfileNum, 0)

	err = dao.ReportError.DB.Model(dao.ReportError.Table).
		Where("report_type = ? and error_type = ? and create_date >= ? and create_date < ?", reportType, errorType, startTime, endTime).
		Fields("count(*) as num, profile_id").
		Group("profile_id").
		Scan(&profiles)
	return profiles, err
}

func (t *ReportErrorRepository) ClearByDate(date time.Time) error {

	for {
		res, err := dao.ReportError.DB.Model(dao.ReportError.Table).
			Limit(500).Delete("create_date < ?", date)

		if err != nil {
			return err
		}

		num, err2 := res.RowsAffected()
		log.Println("error", num)

		if num == 0 || err2 != nil {
			return err2
		}
	}
}



