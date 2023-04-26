package repo

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/google/wire"
	"log"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"time"
)

var ReportSbConsumerDetailRepositorySet = wire.NewSet(wire.Struct(new(ReportSbConsumerDetailRepository), "*"))

type ReportSbConsumerDetailRepository struct{}

func NewReportSbConsumerDetailRepository() *ReportSbConsumerDetailRepository {
	return &ReportSbConsumerDetailRepository{}
}

func (t *ReportSbConsumerDetailRepository) Addone(mqdata *bean.ConsumerDetail) {
	data := model.ReportSbConsumerDetail{
		CtxId:      mqdata.CtxId,
		CreateTime: gtime.NewFromTime(time.Unix(mqdata.CreateTime, 0)),
		ReportName: mqdata.ReportName,
		ProfileId:  mqdata.ProfileId,
		Status:     mqdata.Status,
		ReportDate: mqdata.ReportDate,
		Error:      mqdata.ErrDesc,
		Batch:      mqdata.Batch,
		CostTime:   int(mqdata.CostTime),
		UpdateTime: gtime.NewFromTime(time.Unix(mqdata.UpdateTime, 0)),
	}

	dao.ReportSbConsumerDetail.DB.Model(dao.ReportSbConsumerDetail.Table).Data(data).Insert()
}

func (t *ReportSbConsumerDetailRepository) ClearByDate(date time.Time) error {

	for {
		res, err := dao.ReportSbConsumerDetail.DB.Model(dao.ReportSbConsumerDetail.Table).
			Limit(500).Delete("create_time < ?", date)

		if err != nil {
			return err
		}

		num, err2 := res.RowsAffected()
		log.Println("sb", num)

		if num == 0 || err2 != nil {
			return err2
		}
	}
}
