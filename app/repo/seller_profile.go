package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"td_report/boot"
	"td_report/common/redis"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/google/wire"
)

var SellerProfileRepositorySet = wire.NewSet(wire.Struct(new(SellerProfileRepository), "*"))

type SellerProfileRepository struct{}

func NewSellerProfileRepository() *SellerProfileRepository {
	return &SellerProfileRepository{}
}

func (t *SellerProfileRepository) GetOneByprofiled(profiled int64) (*model.SellerProfile, error) {
	var (
		result        gdb.Record
		SellerProfile *model.SellerProfile
		err           error
	)

	result, err = g.DB().Model(dao.SellerProfile.Table).
		Where("profileId =?", profiled).One()
	if err != nil || result == nil {
		return nil, err
	}

	if result.IsEmpty() {
		return nil, errors.New("no data")
	}

	err = result.Struct(&SellerProfile)
	if err == nil {
		return SellerProfile, nil
	} else {
		return nil, err
	}
}

func (t *SellerProfileRepository) GetOneData(parama map[string]interface{}) (*model.SellerProfile, error) {
	var (
		result        gdb.Record
		SellerProfile *model.SellerProfile
		err           error
	)

	result, err = g.DB().Model(dao.SellerProfile.Table).
		Where("profileId =? and campaignType=? and asin=? and sku=?", parama).One()
	if err != nil || result == nil {
		return nil, err
	}

	if result.IsEmpty() {
		return nil, errors.New("no data")
	}

	err = result.Struct(&SellerProfile)
	if err == nil {
		return SellerProfile, nil
	} else {
		return nil, err
	}
}

func (t *SellerProfileRepository) GetSellerProfileId() ([]string, error) {
	var (
		profileIdList []string
		err           error
		result        []gdb.Record
	)

	result, err = dao.SellerProfile.DB.Model(dao.SellerProfile.Table).Fields("profileId").All()
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

//func (t *SellerProfileRepository) ListProfileAndRefreshtokenFilter(profileIdList []string) ([]*bean.ProfileToken, error) {
//	var sqlStr string
//	if len(profileIdList) != 0 {
//		sqlStr = `
//SELECT
//a.profileId,
//a.region,
//b.refreshToken,
//FROM seller_profile a
//left join (SELECT * FROM seller_token) b
//on a.tokenId = b.id
//where a.status = 0 and a.type = "seller" and b.isPPC = 1 and b.refreshToken != "" %s
//order by a.profileId;
//`
//		str := " and profileId not in ("
//		for _, item := range profileIdList {
//			str = str + item + ","
//		}
//
//		str = strings.Trim(str, ",")
//		str = str + ") "
//		sqlStr = fmt.Sprintf(sqlStr, str)
//
//	} else {
//		sqlStr = `SELECT
//a.profileId,a.region,b.refreshToken
//FROM seller_profile a
//left join (SELECT * FROM seller_token) b
//on a.tokenId = b.id where a.status = 0 and a.type = "seller" and b.isPPC = 1
//and b.refreshToken != "" order by a.profileId;`
//	}
//
//	var token []*bean.ProfileToken
//	err := g.DB().Model(dao.SellerProfile.Table).Raw(sqlStr).Scan(&token)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//
//	return token, nil
//}

func (t *SellerProfileRepository) GetProfileAndRefreshTokenByFilter(profileIdList []string) ([]*bean.ProfileToken, error) {
	data := make([]*bean.ProfileToken, 0)

	query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table + " sp").
		LeftJoin(dao.SellerToken.Table + " st", "sp.tokenId = st.id").
		LeftJoin(dao.SellerClient.Table + " sc", "st.clientId = sc.id").
		Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag, sp.type as profileType").
		Where(" sp.status = ? and st.isPPC = ? and st.refreshToken != ? ", 0, 1, "").
		WhereIn("sp.type", []string{"seller", "vendor"})
		if len(profileIdList) > 0 {
			query.WhereNotIn("sp.profileId", profileIdList)
		}

	err := query.Order("sp.id asc").Scan(&data)
	return data, err
}

func (t *SellerProfileRepository) GetProfileAndRefreshTokenById(profileIdList []string) ([]*bean.ProfileToken, error) {
	data := make([]*bean.ProfileToken, 0)
	query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table+" sp").
		LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
		LeftJoin(dao.SellerClient.Table+" sc", "st.clientId = sc.id").
		Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag, sp.type as profileType").
		Where(" sp.status = ? and st.isPPC = ? and st.refreshToken != ? ", 0, 1, "").
		WhereIn("sp.type", []string{"seller", "vendor"}).
		WhereIn("sp.profileId", profileIdList)
	err := query.Scan(&data)
	return data, err
}

// GetDisableProfile 获取所有ppc禁用用户
func (t *SellerProfileRepository) GetDisableProfile() ([]*bean.ProfileToken, error) {
	data := make([]*bean.ProfileToken, 0)
	query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table + " sp").
		LeftJoin(dao.SellerToken.Table + " st", "sp.tokenId = st.id").
		LeftJoin(dao.SellerClient.Table + " sc", "st.clientId = sc.id").
		Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag, sp.type as profileType").
		Where(" sp.status = ? and st.isPPC = ? and st.refreshToken != ? ", 0, 1, "").
		WhereIn("sp.type", []string{"seller", "vendor"})

	err := query.Scan(&data)
	return data, err
}

func (t *SellerProfileRepository) ListProfileAndRefreshtoken(profileIdList []string) ([]*bean.ProfileToken, error) {
	var err error
	var token []*bean.ProfileToken
	if len(profileIdList) != 0 {
		query := dao.SellerProfile.As("sp").
			LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
			LeftJoin(dao.SellerClient.Table+" sc", "st.clientId = sc.id").
			Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag").
			Where(" sp.status = ? and sp.type = ? and st.isPPC = ? and st.refreshToken != ? ", 0, "seller", 1, "").
			WhereIn("sp.profileId", profileIdList)
		err = query.Scan(&token)

	} else {

		query := dao.SellerProfile.As("sp").
			LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
			LeftJoin(dao.SellerClient.Table+" sc", "st.clientId = sc.id").
			Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag").
			Where(" sp.status = ? and sp.type = ? and st.isPPC = ? and st.refreshToken != ? ", 0, "seller", 1, "")
		err = query.Scan(&token)
	}


	return token, err
}

func (t *SellerProfileRepository) ListProfileAndRefreshtokenMap(profileIdList []string) (map[string]*bean.ProfileToken, error) {
	var (
		err error
		token []*bean.ProfileToken
	)

	if len(profileIdList) != 0 {
		query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table+" sp").
			LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
			LeftJoin(dao.SellerClient.Table+" sc", "st.clientId = sc.id").
			Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag").
			Where(" sp.status = ? and sp.type = ? and st.isPPC = ? and st.refreshToken != ? ", 0, "seller", 1, "").WhereIn("sp.profileId", profileIdList)
		err = query.Scan(&token)

	} else {

		query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table+" sp").
			LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
			LeftJoin(dao.SellerClient.Table+" sc", "st.clientId = sc.id").
			Fields("sp.profileId, sp.region, st.refreshToken, sc.clientId, sc.clientSecret, sc.id as tag").
			Where(" sp.status = ? and sp.type = ? and st.isPPC = ? and st.refreshToken != ? ", 0, "seller", 1, "").WhereIn("sp.profileId", profileIdList)
		err = query.Scan(&token)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(token) != 0 {
		var profileTokenMap = make(map[string]*bean.ProfileToken, 0)
		for _, item := range token {
			profileTokenMap[item.ProfileId] = item
		}

		return profileTokenMap, nil
	} else {
		return nil, errors.New("error data is empty")
	}
}

func (t *SellerProfileRepository) ListAllProfile() ([]*bean.ProfileToken, error) {
	val, err := boot.RedisCommonClient.GetClient().Get(redis.WithProfileTokenKey()).Bytes()
	if err != nil {
		return nil, err
	}

	pts := make([]*bean.ProfileToken, 0)
	return pts, json.Unmarshal(val, &pts)
}

func (t *SellerProfileRepository) GetAllData() (*[]model.SellerProfile, error) {
	var (
		SellerProfile *[]model.SellerProfile
		err           error
	)

	err = g.DB().Model(dao.SellerProfile.Table).Where("status!=0").Scan(&SellerProfile)
	if err != nil {
		return nil, err
	}

	return SellerProfile, nil
}

// GetProfile datatype 1 普通 2 sb的特殊报表
func (t *SellerProfileRepository) GetProfile(datatype int) ([]*bean.ProfileToken, error) {
	var sqlStr string
	if datatype == 1 {
		sqlStr = `
SELECT a.profileId, a.region, b.refreshToken
FROM seller_profile a left join (SELECT * FROM seller_token) b
on a.tokenId = b.id where a.status = 0 and a.type = "seller" and b.isPPC = 1 
and b.refreshToken != "" order by a.profileId;`

	} else {
		sqlStr = `
SELECT distinct a.profileId,a.region,b.refreshToken FROM seller_profile a
inner join brand t on t.profileId = a.profileId  
left join (SELECT * FROM seller_token) b on a.tokenId = b.id
where a.status = 0 and a.type = "seller" and b.isPPC = 1 and b.refreshToken != ""
order by a.profileId;`
	}

	var (
		SellerProfile []*bean.ProfileToken
		err           error
	)

	err = g.DB().Raw(sqlStr).Scan(&SellerProfile)
	if err != nil {
		return nil, err
	}

	return SellerProfile, nil
}

func (t *SellerProfileRepository) ArrayInGroupsOf(arr []*bean.ProfileToken, num int64) [][]*bean.ProfileToken {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]*bean.ProfileToken{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]*bean.ProfileToken, 0)
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

func (t *SellerProfileRepository) ArrayInGroupsOfString(arr []string, num int64) [][]string {
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

func (t *SellerProfileRepository) GetProfileAndRefreshToken(clientTag int, profileIdList ...int64) ([]*bean.ProfileTokenClient, error) {
	data := make([]*bean.ProfileTokenClient, 0)
	query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table+" sp").
		LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
		LeftJoin(dao.SellerClient.Table+" sc", "st.clientId = sc.id").
		Fields("sp.profileId, sp.region, st.refreshToken, st.clientId as client_tag , sc.clientId, sc.clientSecret").
		Where(" sp.status = ? and sp.type = ? and st.isPPC = ? and st.refreshToken != ? ", 0, "seller", 1, "")
	if len(profileIdList) > 0 {
		query.WhereIn("sp.profileId", profileIdList)
	}

	if clientTag != 0 {
		query.Where("st.clientId=?", clientTag)
	}

	err := query.Scan(&data)
	return data, err
}

func (t *SellerProfileRepository) GetAll() (map[string]string, error) {
	data := make([]*bean.ProfileAll, 0)
	query := dao.SellerProfile.DB.Model(dao.SellerProfile.Table+" sp").
		LeftJoin(dao.SellerToken.Table+" st", "sp.tokenId = st.id").
		Fields("sp.profileId, sp.region, st.nickname").
		Where(" sp.status = ? and sp.type = ? and st.isPPC = ? and st.refreshToken != ? ", 0, "seller", 1, "")

	err := query.Scan(&data)
	if err != nil {
		return nil, err
	}
	datamap := make(map[string]string, 0)
	for _, item := range data {
		key := fmt.Sprintf("%d", item.ProfileId)
		datamap[key] = item.NickName
	}
	return datamap, nil
}
