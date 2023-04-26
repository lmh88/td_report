package get_token

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/common/tool"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"

	"td_report/pkg/requests"
	"time"
)


func GetAccessToken(crt *varible.ClientRefreshToken) (string, error) {
	token, err := getAccessTokenCache(crt.RefreshToken)
	if err != nil && err != redis.Nil {
		return "", err
	}

	if err == redis.Nil || token == "" {

		defer boot.Redisclient.SyncUnlock(crt.RefreshToken)
		if boot.Redisclient.SyncLock(crt.RefreshToken) {
			rToken, resp, err := RefreshToken(crt)
			if err != nil {
				return "", err
			}

			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return "", errors.New(fmt.Sprintf("%d,%s", resp.StatusCode, resp.Body))
			}
			key := tool.GetMd5(crt.RefreshToken)
			err = boot.Redisclient.GetClient().Set(redis.WithAccessTokenPrefix(key),
				rToken.AccessToken,
				time.Second * time.Duration(rToken.ExpiresIn-10)).Err()
			if err != nil {
				return "", err
			}
			return rToken.AccessToken, nil
		} else {
			for i := 0; i < 3; i++ {
				time.Sleep(time.Second * 3)
				token, err = getAccessTokenCache(crt.RefreshToken)
				if err != nil || token == "" {
					if i < 2 {
						continue
					}
				}
				return token, err
			}
		}
	}

	return token, nil
}

func getAccessTokenCache(refreshToken string) (string, error) {
	key := tool.GetMd5(refreshToken)
	return boot.Redisclient.GetClient().Get(redis.WithAccessTokenPrefix(key)).Result()
}

type RTokenBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func RefreshToken(crt *varible.ClientRefreshToken) (*RTokenBody, *requests.Resp,  error) {

	url := "https://api.amazon.com/auth/o2/token"
	var limit = limiter.PerSecond(g.Cfg().GetInt("common.token_limit"))
	//添加频率率的限制，避免请求频次太多
	for {

		res, err := boot.Rlimiter.Allow(url, limit)
		if err != nil {
			return nil, nil, err
		}

		if res.Allowed == 1 {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	params := map[string]interface{}{
		"grant_type":    "refresh_token",
		"client_id":     crt.ClientId,
		"refresh_token": crt.RefreshToken,
		"client_secret": crt.ClientSecret,
	}
	logger.Info(url, params)

	var rToken RTokenBody
	i := 0
	for {
		resp, err := requests.Post(url, requests.WithJson(params), requests.WithTimeout(time.Second*60))
		if err != nil {
			logger.Logger.Error(err, " post method get token error")
			i++
			if i >= 3 {
				return nil, resp, err
			}
			time.Sleep(time.Second)
			continue
		}

		_, err = resp.Json(&rToken)
		if err != nil {
			logger.Logger.Error(err, " json get token error ", string(resp.Body))
			return nil, resp, err
		}
		if rToken.ExpiresIn == 0 {
			i++
			if i < 3 {
				time.Sleep(time.Second)
				continue
			}
		}
		return &rToken, resp, nil
	}
}

//func getUrlTest(url string) (*requests.Resp,  error) {
//	var (
//		resp requests.Resp
//		err error
//	)
//
//	fmt.Println(url)
//	url = "{\"expires_in\":3600, \"access_token\":\"xyz\"}"
//	resp = requests.Resp{
//		StatusCode: 200,
//		Body: []byte(url),
//	}
//
//	//err = errors.New("chang error")
//	err = nil
//
//	return &resp, err
//}
