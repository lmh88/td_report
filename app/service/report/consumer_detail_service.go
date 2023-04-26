package report

import (
	"github.com/google/wire"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/vars"
)

var ConsumerDetailServiceSet = wire.NewSet(wire.Struct(new(ConsumerDetailService), "*"))

type ConsumerDetailService struct {
	ConsumerDetailSdRepository  *repo.ReportSdConsumerDetailRepository
	ConsumerDetailSpRepository  *repo.ReportSpConsumerDetailRepository
	ConsumerDetailSbRepository  *repo.ReportSbConsumerDetailRepository
	ConsumerDetailDspRepository *repo.ReportDspConsumerDetailRepository
}

func NewConsumerDetailService() *ConsumerDetailService {
	return &ConsumerDetailService{
		ConsumerDetailSbRepository:  repo.NewReportSbConsumerDetailRepository(),
		ConsumerDetailSdRepository:  repo.NewReportSdConsumerDetailRepository(),
		ConsumerDetailSpRepository:  repo.NewReportSpConsumerDetailRepository(),
		ConsumerDetailDspRepository: repo.NewReportDspConsumerDetailRepository(),
	}
}

func (t *ConsumerDetailService) AddConsumerDetail(mqdata *bean.ConsumerDetail) {
	switch mqdata.ReportType {
	case vars.SP:
		t.ConsumerDetailSpRepository.Addone(mqdata)
	case vars.SB:
		t.ConsumerDetailSbRepository.Addone(mqdata)
	case vars.SD:
		t.ConsumerDetailSdRepository.Addone(mqdata)
	case vars.DSP:
		t.ConsumerDetailDspRepository.Addone(mqdata)
	}
}
