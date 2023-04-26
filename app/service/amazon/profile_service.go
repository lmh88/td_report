package amazon

import (
	"encoding/json"
	"github.com/google/wire"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/pkg/logger"
	"time"
)

// ProfileService 设置dsp类型的区域信息
type ProfileService struct {
	ProfileRepository       *repo.ProfileRepository
	SellerProfileRepository *repo.SellerProfileRepository
}

var ProfileServiceSet = wire.NewSet(wire.Struct(new(ProfileService), "*"))

func NewProfileService(profileRepository *repo.ProfileRepository) *ProfileService {
	return &ProfileService{
		ProfileRepository: profileRepository,
	}
}

func (t *ProfileService) GetPpcProfile(profileIdList []string) ([]*bean.ProfileToken, error) {
	return t.SellerProfileRepository.ListProfileAndRefreshtoken(profileIdList)
}

func (t *ProfileService) GetPpcProfileFilter(profileIdList []string) ([]*bean.ProfileToken, error) {
	return t.SellerProfileRepository.GetProfileAndRefreshTokenByFilter(profileIdList)
}

func (t *ProfileService) GetPpcProfileMap(profileIdList []string) (map[string]*bean.ProfileToken, error) {
	return t.SellerProfileRepository.ListProfileAndRefreshtokenMap(profileIdList)
}

// GetProfile profileType =1 ppc 2 dsp
func (t *ProfileService) GetProfile(profileType int) ([]string, error) {
	if profileType == 1 {
		return t.SellerProfileRepository.GetSellerProfileId()
	}

	return t.ProfileRepository.GetDspProfileId()
}

func (t *ProfileService) SetPpcProfile() {
	profileIdList := make([]string, 0)
	pts, err := t.GetPpcProfile(profileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListProfileAndRefreshtoken error",
			"err":  err.Error(),
		})
		return
	}

	val, err := json.Marshal(pts)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListProfileAndRefreshtoken error",
			"err":  err.Error(),
		})
		return
	}

	err = boot.RedisCommonClient.GetClient().Set(redis.WithProfileTokenKey(), val, time.Hour*3).Err()
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListProfileAndRefreshtoken error",
			"err":  err.Error(),
		})
		return
	}
}

func (t *ProfileService) GetDspProfile(profileIdList []string) ([]*bean.DspRegionProfile, error) {
	return t.ProfileRepository.ListDspRegionProfile(profileIdList)
}

func (t *ProfileService) GetDspProfileFilter(profileIdList []string) ([]*bean.DspRegionProfile, error) {
	return t.ProfileRepository.ListDspRegionProfileFilter(profileIdList)
}

func (t *ProfileService) SetDspregin() {
	var profileIdList = make([]string, 0)
	rps, err := t.ProfileRepository.ListDspRegionProfile(profileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListDspRegionProfile error",
			"err":  err.Error(),
		})
		return
	}

	val, err := json.Marshal(rps)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListDspRegionProfile error",
			"err":  err.Error(),
		})
		return
	}

	err = boot.RedisCommonClient.GetClient().Set(redis.WithDspRegionProfile(), val, time.Hour*3).Err()
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListDspRegionProfile error",
			"err":  err.Error(),
		})
		return
	}
}
