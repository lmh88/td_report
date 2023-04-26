package repo

import (
	"github.com/gogf/gf/database/gdb"
	"td_report/app/dao"
)

type SellerClientRepository struct{}

func NewSellerClientRepository() *SellerClientRepository {
	return &SellerClientRepository{}
}

func (t *SellerClientRepository) GetOne(condition map[string]int64) (gdb.Record, error) {
	return dao.SellerProfile.DB.Model(dao.SellerProfile.Table).Where(condition).One()
}

func (t *SellerClientRepository) GetAll() []int {
	data, _ := dao.SellerClient.DB.Model(dao.SellerClient.Table).Fields("id").All()
	res := make([]int, 0)
	for _, item := range data {
		res = append(res, item.Map()["id"].(int))
	}
	return res
}
