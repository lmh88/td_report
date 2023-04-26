package dbp_token

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/common/tool"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"
	"td_report/pkg/requests"
	"time"
)

func GetAccessTokenWithClient(refreshToken, clientId, clientSecret string) (string, error) {
	token, err := GetAccessToken(refreshToken)
	if err != nil && err != redis.Nil {
		return "", err
	}

	if err == redis.Nil || token == "" {
		return getAccessTokenAndRefreshTokenWithClient(refreshToken, clientId, clientSecret)
	}

	return token, nil
}

func GetAccessToken(refreshToken string) (string, error) {
	key := tool.GetMd5(refreshToken)
	return boot.Redisclient.GetClient().Get(redis.WithAccessTokenPrefix(key)).Result()
}

func getAccessTokenAndRefreshTokenWithClient(refreshToken, clientId, clientSecret string) (string, error) {
	defer boot.Redisclient.SyncUnlock(refreshToken)
	if boot.Redisclient.SyncLock(refreshToken) {
		var (
			n      = 20
			rToken *RTokenBody
			err    error
			code   int
			reason string
		)

		for i := 0; i < n; i++ {
			rToken, code, reason, err = getRefreshTokenWithClient(refreshToken, clientId, clientSecret)
			if err != nil {
				return "", err
			}
			//发现时间为0的情况循环调用30次
			if rToken.ExpiresIn != 0 {
				break
			} else {

				// 400 状态
				if code == 400 {
					return "", errors.New(reason)
				}
				time.Sleep(5 * time.Second)
				logger.Logger.Info(fmt.Sprintf("token expire time is 0, try times%d and sleep 1s", i))
			}
		}

		key := tool.GetMd5(refreshToken)
		err = boot.Redisclient.GetClient().Set(
			redis.WithAccessTokenPrefix(key),
			rToken.AccessToken,
			time.Second*time.Duration(rToken.ExpiresIn-50)).Err()
		if err != nil {
			return "", err
		}

		return rToken.AccessToken, nil

	} else {

		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 3)
			token, err := GetAccessToken(refreshToken)
			if err != nil || token == "" {
				if i < 2 {
					continue
				}
			}
			return token, err
		}
	}

	return "", errors.New("error")
}

func getRefreshTokenWithClient(refreshToken, clientId, clientSecret string) (*RTokenBody, int, string, error) {
	var rlimit = limiter.PerSecond(g.Cfg().GetInt("common.token_limit"))
	url := "https://api.amazon.com/auth/o2/token"
	//添加频率率的限制，避免请求频次太多
	for {

		res, err := boot.Rlimiter.Allow(url, rlimit)
		if err != nil {
			return nil, 200, "", err
		}

		if res.Allowed == 1 {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	params := map[string]interface{}{
		"grant_type":    "refresh_token",
		"client_id":     clientId,
		"refresh_token": refreshToken,
		"client_secret": clientSecret,
	}

	var rToken RTokenBody
	resp, err := requests.Post(url, requests.WithJson(params), requests.WithTimeout(time.Second*60))
	if err != nil {
		logger.Logger.Error(err, "post method get  token error")
		return nil, 200, "post method get  token error", err
	}

	_, err = resp.Json(&rToken)
	if err != nil {
		logger.Logger.Error(err, "get  token error")
		return nil, 200, "trans token json  error", err
	}

	if rToken.ExpiresIn == 0 {
		logger.Logger.Info("token ========", resp.StatusCode, string(resp.Body))
	}

	return &rToken, resp.StatusCode, string(resp.Body), nil
}
