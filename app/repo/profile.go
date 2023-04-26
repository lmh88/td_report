package repo

import (
	"encoding/json"
	"errors"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/google/wire"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"td_report/boot"
	"td_report/common/redis"
)

var ProfileRepositorySet = wire.NewSet(wire.Struct(new(ProfileRepository), "*"))

type ProfileRepository struct{}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{}
}

func (t *ProfileRepository) GetOneByprofiled(profiled int64) (*model.Profile, error) {
	var (
		result  gdb.Record
		Profile *model.Profile
		err     error
	)

	result, err = g.DB().Model(dao.Profile.Table).
		Where("profileId =?", profiled).One()
	if err != nil || result == nil {
		return nil, err
	}

	if result.IsEmpty() == true {
		return nil, errors.New("no data")
	}

	err = result.Struct(&Profile)
	if err == nil {
		return Profile, nil
	} else {
		return nil, err
	}
}

func (t *ProfileRepository) ListDspProfile() ([]*bean.DspRegionProfile, error) {
	val, err := boot.RedisCommonClient.GetClient().Get(redis.WithDspRegionProfile()).Bytes()
	if err != nil {
		return nil, err
	}

	rps := make([]*bean.DspRegionProfile, 0)
	return rps, json.Unmarshal(val, &rps)
}

func (t *ProfileRepository) ListDspRegionProfile(profileIdList []string) ([]*bean.DspRegionProfile, error) {
	var (
		Profiles []*bean.DspRegionProfile
		err      error
	)

	model := dao.Profile.DB.Model(dao.Profile.Table).Fields("region, profileId ").
		Where("type =? and status=0", "agency")
	if len(profileIdList) > 0 {
		model.WhereIn("profileId", profileIdList)
	}

	err = model.Scan(&Profiles)
	if err != nil {
		return nil, err
	}
	return Profiles, err
}

func (t *ProfileRepository) ListDspRegionProfileFilter(profileIdList []string) ([]*bean.DspRegionProfile, error) {
	var (
		Profiles []*bean.DspRegionProfile
		err      error
	)

	model := dao.Profile.DB.Model(dao.Profile.Table).Fields("region, profileId ").
		Where("type =? and status=0", "agency")
	if len(profileIdList) > 0 {
		model.WhereNotIn("profileId", profileIdList)
	}

	err = model.Scan(&Profiles)
	if err != nil {
		return nil, err
	}
	return Profiles, err
}

func (t *ProfileRepository) GetDspProfileId() ([]string, error) {
	var (
		profileIdList []string = make([]string, 0)
		err           error
		result        []gdb.Record
	)

	result, err = dao.Profile.DB.Model(dao.Profile.Table).Fields("profileId").All()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return profileIdList, nil
	}
	for _, item := range result {
		profileIdList = append(profileIdList, item["profileId"].String())
	}
	return profileIdList, err
}

func (t *ProfileRepository) ArrayInGroupsOf(arr []*bean.DspRegionProfile, num int64) [][]*bean.DspRegionProfile {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]*bean.DspRegionProfile{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]*bean.DspRegionProfile, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}

	return segments
}

func (t *ProfileRepository) ArrayInGroupsOfString(arr []string, num int64) [][]string {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]string{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]string, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}

	return segments
}
