package repo

import (
	"github.com/google/wire"
	"td_report/app/dao"
	"td_report/app/model"
)

var SpTrafficRepositorySet = wire.NewSet(wire.Struct(new(SpTrafficRepository), "*"))

type SpTrafficRepository struct{}

func NewSpTrafficRepository() *SpTrafficRepository {
	return &SpTrafficRepository{}
}

// AddMutils 单个或者多个插入数据
func (t *SpTrafficRepository) AddMutils(data model.SpTraffic) error {
	_, err := dao.SpTraffic.Data(data).Insert()
	return err
}
