package common

import (
	"errors"
	"github.com/google/wire"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/vars"
)

type AuthService struct {
	SellerProfileRepo *repo.SellerProfileRepository
	SellerTokenRepo   *repo.SellerTokenRepository
	ProfileRepo       *repo.ProfileRepository
}

var AuthServiceSet = wire.NewSet(wire.Struct(new(AuthService), "*"))

func NewAuthService(SellerProfileRepository *repo.SellerProfileRepository,
	SellerTokenRepository *repo.SellerTokenRepository, ProfileRepository *repo.ProfileRepository) *AuthService {
	return &AuthService{
		SellerProfileRepo: SellerProfileRepository,
		SellerTokenRepo:   SellerTokenRepository,
		ProfileRepo:       ProfileRepository,
	}
}

// GetAuth 针对sp sd sb获取权限参数
func (t *AuthService) GetAuth(profileId int64, datatype string) (map[string]interface{}, error) {
	var (
		sellerProfile *model.SellerProfile
		err           error
	)
	sellerProfile, err = t.SellerProfileRepo.GetOneByprofiled(profileId)
	if err != nil {
		return nil, err
	}
	if sellerProfile == nil {
		return nil, errors.New("授权失败profiled找不到数据")
	}

	authdata := make(map[string]interface{}, 0)

	switch datatype {
	case vars.SP:
		var paramas vars.SpParams
		paramas.BaseParams.TokenId = sellerProfile.TokenId
		paramas.BaseParams.ProfileId = sellerProfile.ProfileId
		paramas.BaseParams.Region = vars.RegionType(sellerProfile.Region)
		authdata[vars.SP] = paramas
	case vars.SD:
		var paramas vars.SdParams
		paramas.BaseParams.TokenId = sellerProfile.TokenId
		paramas.BaseParams.ProfileId = sellerProfile.ProfileId
		paramas.BaseParams.Region = vars.RegionType(sellerProfile.Region)
		authdata[vars.SD] = paramas
	case vars.SB:
		var paramas vars.SbParams
		paramas.BaseParams.TokenId = sellerProfile.TokenId
		paramas.BaseParams.ProfileId = sellerProfile.ProfileId
		paramas.BaseParams.Region = vars.RegionType(sellerProfile.Region)
		authdata[vars.SB] = paramas
	default:
		return nil, errors.New(" 非法参数")
	}

	return authdata, nil
}

// GetAuthDsp 针对dsp
func (t *AuthService) GetAuthDsp(profileId int64) (map[string]interface{}, error) {
	var (
		Profile *model.Profile
		err     error
	)

	Profile, err = t.ProfileRepo.GetOneByprofiled(profileId)
	if err != nil {
		return nil, err
	}

	authdata := make(map[string]interface{}, 0)
	authdata["regin"] = vars.RegionType(Profile.Region)
	authdata["profileId"] = Profile.ProfileId
	authdata["tokenId"] = 0
	return authdata, nil
}
