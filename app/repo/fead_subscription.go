package repo

import (
	"github.com/gogf/gf/frame/g"
	"td_report/app/dao"
	"td_report/app/model"
)

type FeadSubscriptionRepository struct {}

func NewFeadSubscriptionRepository() *FeadSubscriptionRepository {
	return new(FeadSubscriptionRepository)
}

func (r *FeadSubscriptionRepository) AddBatch(items []*model.FeadSubscription) error {
	_, err := dao.FeadSubscription.DB().Model(dao.FeadSubscription.Table).Data(items).Save()
	//fmt.Println(res.RowsAffected())
	//fmt.Println(res.LastInsertId())
	return err
}

func (r *FeadSubscriptionRepository) UpdateOne(id uint, data g.Map) error {
	_, err := dao.FeadSubscription.DB().Model(dao.FeadSubscription.Table).Data(data).Where("id", id).Update()
	return err
}

func (r *FeadSubscriptionRepository) GetAll() []*model.FeadSubscription {
	res := make([]*model.FeadSubscription, 0)
	dao.FeadSubscription.DB().Model(dao.FeadSubscription.Table).
		Where("is_del", 0).Scan(&res)
	return res
}

func (r *FeadSubscriptionRepository) GetByProfile(profileId string) []*model.FeadSubscription {
	res := make([]*model.FeadSubscription, 0)
	dao.FeadSubscription.DB().Model(dao.FeadSubscription.Table).
		Where("is_del = ? and profile_id = ?", 0, profileId).
		Scan(&res)
	return res
}

func (r *FeadSubscriptionRepository) DelOne(id uint) error {
	_, err := dao.FeadSubscription.DB().Model(dao.FeadSubscription.Table).Data(g.Map{"is_del": 1}).Where("id", id).Update()
	return err
}

func (r *FeadSubscriptionRepository) GetBySubId(subId string) *model.FeadSubscription {
	//res := make([]*model.FeadSubscription, 0)
	var one model.FeadSubscription
	dao.FeadSubscription.DB().Model(dao.FeadSubscription.Table).
		Where("is_del = ? and message_subscription_id = ?", 0, subId).
		Scan(&one)
	return &one
}
