package repo

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/google/wire"
	"td_report/app/dao"
	"td_report/app/model"
)

var FeadRepositorySet = wire.NewSet(wire.Struct(new(FeadRepository), "*"))

type FeadRepository struct{}

func NewFeadRepository() *FeadRepository {
	return &FeadRepository{}
}

func (t *FeadRepository) AddOne(data *model.Fead) error {
	var da *model.Fead
	err := dao.Fead.Where("profile_id=? and dataset_id=?", data.ProfileId, data.DatasetId).Scan(&da)
	if err != nil {
		return err
	}
	if da != nil {
		data.UpdateDate = gtime.Now()
		data.Id = da.Id
		if data.MessagesSubscriptionId == "" {
			data.MessagesSubscriptionId = da.MessagesSubscriptionId
		}

		if data.ClientTokenId == "" {
			data.ClientTokenId = da.ClientTokenId
		}

		if data.Version == 0 {
			data.Version = da.Version
		}

		data.UpdateDate = gtime.Now()
		_, err = dao.Fead.Data(data).Where("profile_id=? and dataset_id=?", data.ProfileId, data.DatasetId).Update()
		fmt.Println(data.ProfileId, " 添加fead中，已有历史数据，修改数据")
		return err
	} else {
		_, err = dao.Fead.Data(data).Insert()
		return err
	}
}

func (t *FeadRepository) Update(profileId int64, datdaSetId string, paramas map[string]interface{}) error {
	paramas["update_date"] = gtime.Now()
	_, err := dao.Fead.Data(paramas).Where("profile_id=? and dataset_id=?", profileId, datdaSetId).Update()
	return err
}
