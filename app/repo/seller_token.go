package repo

import (
	"errors"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/google/wire"
	"td_report/app/dao"
	"td_report/app/model"
)

type SellerTokenRepository struct{}

var SellerTokenRepositorySet = wire.NewSet(wire.Struct(new(SellerTokenRepository), "*"))

func NewSellerTokenRepository() *SellerTokenRepository {
	return &SellerTokenRepository{}
}

func (t *SellerProfileRepository) GetOne(condition map[string]int64) (gdb.Record, error) {
	return dao.SellerProfile.DB.Model(dao.SellerProfile.Table).Where(condition).One()
}

func (t *SellerTokenRepository) GetOneData(tokenid int) (*model.SellerToken, error) {
	var (
		result      gdb.Record
		SellerToken *model.SellerToken
		err         error
	)

	result, err = g.DB().Model(dao.SellerToken.Table).Where("id=?", tokenid).One()
	if err != nil || result == nil {
		return nil, err
	}

	if result.IsEmpty() == true {
		return nil, errors.New("no data")
	}

	err = result.Struct(&SellerToken)
	if err == nil {
		return SellerToken, nil
	} else {
		return nil, err
	}
}
