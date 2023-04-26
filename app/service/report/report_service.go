package report

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
	"td_report/app/bean"
	"td_report/app/repo"
)

var AllReportServiceSet = wire.NewSet(wire.Struct(new(AllReportService), "*"))

type AllReportService struct {
	SellerProfileRepository *repo.SellerProfileRepository
	ProfileRepository       *repo.ProfileRepository
	ReportTaskRepository    *repo.ReportTaskRepository
}

func NewAllReportService(SellerProfileRepository *repo.SellerProfileRepository,
	ProfileRepository *repo.ProfileRepository,
	ReportTaskRepository *repo.ReportTaskRepository,
) *AllReportService {
	return &AllReportService{
		SellerProfileRepository: SellerProfileRepository,
		ProfileRepository:       ProfileRepository,
		ReportTaskRepository:    ReportTaskRepository,
	}
}

func (t *AllReportService) ReceiveSp(msg *amqp.Delivery) error {
	g.Log().Println("get queue sp common:", string(msg.Body))
	return nil
}

func (t *AllReportService) ReceiveSd(msg *amqp.Delivery) error {
	g.Log().Println("get queue sd common:", string(msg.Body))
	return nil
}

func (t *AllReportService) ReceiveSb(msg *amqp.Delivery) error {
	g.Log().Println("get queue sb common:", string(msg.Body))
	return nil
}

func (t *AllReportService) ReceiveDsp(msg *amqp.Delivery) error {
	g.Log().Println("get queue dsp common:", string(msg.Body))
	return nil
}

func (t *AllReportService) receiveCommon(msg *amqp.Delivery) (*bean.ReportCommonData, error) {
	var paramas bean.ReceiveMsgBody
	var report *bean.ReportCommonData
	if err := json.Unmarshal(msg.Body, &paramas); err != nil {
		g.Log().Error(err)
		return nil, err
	}

	json.Unmarshal([]byte(paramas.MessageBody), &report)
	return report, nil
}
