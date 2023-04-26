package dsp_region_profile

import (
	"encoding/json"
	"td_report/boot"
	"td_report/common/redis"
)

type DspRegionProfile struct {
	ProfileId string `json:"profile_id"`
	Region    string `json:"region"`
}

func ListDspProfile() ([]*DspRegionProfile, error) {
	val, err := boot.RedisCommonClient.GetClient().Get(redis.WithDspRegionProfile()).Bytes()
	if err != nil {
		return nil, err
	}

	rps := make([]*DspRegionProfile, 0)
	return rps, json.Unmarshal(val, &rps)
}

//func ListDspRegionProfile() ([]*DspRegionProfile, error) {
//	sqlStr := `
//SELECT
//region,
//profileId
//FROM profile
//where type = 'agency' and status = 0
//;
//`
//	rows, err := mysql.DBQuery.Query(sqlStr)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	rps := make([]*DspRegionProfile, 0)
//	for rows.Next() {
//		rp := &DspRegionProfile{}
//		err = rows.Scan(&rp.Region, &rp.ProfileId)
//		if err != nil {
//			return nil, err
//		}
//
//		rps = append(rps, rp)
//	}
//
//	return rps, nil
//}
