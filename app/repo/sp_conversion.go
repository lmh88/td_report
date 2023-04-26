package repo

import (
	"github.com/google/wire"
	"td_report/app/dao"
	"td_report/app/model"
)

var SpConversionRepositorySet = wire.NewSet(wire.Struct(new(SpConversionRepository), "*"))

type SpConversionRepository struct{}

func NewSpConversionRepository() *SpConversionRepository {
	return &SpConversionRepository{}
}

// AddMutils 单个或者多个插入数据
func (t *SpConversionRepository) AddMutils(data model.SpConversion) error {
	_, err := dao.SpConversion.Data(data).Insert()
	return err
}
