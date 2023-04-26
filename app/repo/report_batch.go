package repo

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/google/wire"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
)

var ReportBatchRepositorySet = wire.NewSet(wire.Struct(new(ReportBatchRepository), "*"))

type ReportBatchRepository struct{}

func NewReportBatchRepository() *ReportBatchRepository {
	return &ReportBatchRepository{}
}

func (t *ReportBatchRepository) GetOne(batch string) (*model.ReportBatch, error) {
	var (
		one *model.ReportBatch
		err error
	)
	query, _ := dao.ReportBatch.DB.Model(dao.ReportBatch.Table).Where("batch", batch).One()
	err = query.Struct(&one)
	return one, err
}

func (t *ReportBatchRepository) Addone(uuid string, paramas string) error {
	g := model.ReportBatch{
		Batch:      uuid,
		Paramas:    paramas,
		Status:     1,
		CreateTime: gtime.Now(),
		UpdateTime: gtime.Now(),
		IsCheck:    0,
	}
	_, err := dao.ReportBatch.DB.Model(dao.ReportBatch.Table).Data(g).Insert()
	return err
}

func (t *ReportBatchRepository) ChangeStatus(status int, batch string) error {
	_, err := dao.ReportBatch.DB.Model(dao.ReportBatch.Table).Data(g.Map{"status": status, "update_time": gtime.Now()}).Where("batch", batch).Update()
	return err
}

func (t *ReportBatchRepository) ChangeCheck(isCheck int, batch string) error {
	_, err := dao.ReportBatch.DB.Model(dao.ReportBatch.Table).Data(g.Map{"is_check": isCheck, "update_time": gtime.Now()}).Where("batch", batch).Update()
	return err
}

func (t *ReportBatchRepository) GetList(condition map[string]interface{}) ([]*bean.ReportData, error) {
	var dataList []*model.ReportBatch
	var reportList []*bean.ReportData = make([]*bean.ReportData, 0)
	err := dao.ReportBatch.DB.Model(dao.ReportBatch.Table).Where(condition).OrderAsc("create_time").Scan(&dataList)
	if err != nil {
		return nil, err
	}

	for _, item := range dataList {
		var temp *bean.ReportData
		if item.Paramas != "" {
			if err = json.Unmarshal([]byte(item.Paramas), &temp); err == nil {
				temp.Batch = item.Batch
				reportList = append(reportList, temp)
			}
		}
	}

	return reportList, nil
}

// GetListData 针对队列用的json数据
func (t *ReportBatchRepository) GetListData(condition map[string]interface{}) ([]*model.ReportBatch, error) {
	var dataList []*model.ReportBatch
	err := dao.ReportBatch.DB.Model(dao.ReportBatch.Table).Where(condition).OrderAsc("create_time").Scan(&dataList)
	if err != nil {
		return nil, err
	}

	return dataList, nil
}
