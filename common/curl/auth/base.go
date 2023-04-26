package auth

import (
	"errors"
	"fmt"
	"strconv"
	"td_report/app"
	"td_report/boot"
	"td_report/common/curl"
	"td_report/common/tool"

	"github.com/gogf/gf/frame/g"
)

type RegionType string

type AmazonAuth struct {
	Region    string
	ProfileId string
	TokenId   string
}

const (
	RegionNa RegionType = "NA"
	RegionEu            = "EU"
	RegionFe            = "FE"

	ClientId     string = "amzn1.application-oa2-client.084663234c2143c3a3bf91fe34bbdf1e"
	ClientSecret        = "29f982cfd2571585c52db5ba462f302716c87fd52507dcc355582b8821a80abc"
	ApiTokenUrl         = "https://api.amazon.com/auth/o2/token"
	Version             = "v2"
)

var ApiUrl = map[RegionType]string{
	RegionFe: "https://advertising-api-fe.amazon.com",
	RegionEu: "https://advertising-api-eu.amazon.com",
	RegionNa: "https://advertising-api.amazon.com",
}

var ResToken = map[RegionType]string{
	RegionFe: "Atzr|IwEBIJqR9kHQts1y3pFLrTeLMA8afkEd_dNUAwsjSjqIhFEZHCmVGyhiOkSqZ1gNHUY2X1TSLJjhUeRm6PVe1FnzBAo7o6bOmE495OjbGKaOQAJ2fHh1rBS51dJLT4oX1l7-bWJYb2npLXmOjZQgY8t1ZCDR-45KY0SqzjrbtQMDjNh8uAKRiwhBwBu3H13FcJxxOBi8yjrS2zKuNb2wmWBTe9vUhPbigbvvAMzMeTFmMSK-_N5hd8ZglvJjkvobK5SD7BQeDXK49t5BN-05WSzR6kLfFa9gWSSYphoPdF9hUYaUHIDBUj4ZOeWURoFugZzDO408Xc0lee6GGz1CP21aumiWlSCTuEUrOAENGcNaipEEvWRvZDiHKhIX546ycrVt95ujrQgIlVm5H2-copN3Pkb1KAUBSmW7Ylgwha6OB5fJNg",
	RegionEu: "Atzr|IwEBIJ9_AwBc6Pzv65zeD4ph1nrxNqceokx8IrqxRktMCrpeuC0ZdlqGqUjD8tO3TUCoD8804QB4ND-1Ng-9bY5gMbay9O69vgYN4c2A36y-ep_jzRF4kjyuU22ALAnZjjj6FQ0SoRTkstAZUd7zB3z7U2_i6SIRyaXpDhhIPSAleYEBlQAyM48Lph_VW5dU8pNL7OQrh9r1SW05yyuB7lYaA1juQDYdX4U8WRez0iLssUXhTbCSqTUSSFuqBedAzJMcge84iI_EYp6C9wtFFDFJasXmEfOeLA4S2x_OS07PW06nvXjL3EeLDxgCZqGRTHSVj3V0BbsypQdjYgD7chXeZpAs-YqJKE1dZQ4lPt6yomeJzFFmtO_3RNO1H7_II-nKW3OSHC3WjCG5V9on7_vR1sEI4PZmD0d1RUEEfOL2liTp5ThjRq809SnV9YZL38g9As8",
	RegionNa: "Atzr|IwEBIFohGsLBoxln39hOsfSZfYN18O1YJINunLjrMwBv28y6rpH7EtK8PC19orgSR6TzQE5IkWVOB-e-QMIDCYvDcTCxg9Zqndo9pHLovbPLmXYPdt04DRyEZKgptREAhq8GmZvkZww8ldYm4nGnpFOpV7Czoy8CAjmZHQF9B-jzQhOzjxTKhCUeM2M-6kojkHLmJQqDs3n8QiZr4RVscL3Y41ZLvYwh8D5H1Z93gQNuIdsFzBq0bookPCWJuJeQjC484654ox5AMLfaKlQEejUZQNDdB-Dn6zuvEGudeDPbFh69cANaA7EVHUd4XSl6SP3l0QrW87qLCWz5iXnEnDATVa4-MyGCiMVE-zt4zxxkVJXyDo8UZ8zaFS3qAUdlHOJZ_nGHFHJN7XreJ-cw3foIRLFR1Q6TBYbt6nnmdhFwgj4K2OcePo-Z2T_tvDSIAuLdYn4",
}

// Regions 地区
var Regions = []RegionType{RegionNa, RegionEu, RegionFe}

// GetAccessToken 获取accessToken
func GetAccessToken(region RegionType, tokenId int) string {
	if tokenId == 136 { // 过滤demo的rpc请求
		return ""
	}

	tokenKey := fmt.Sprintf("amazonAccessToken:%s:%d", region, tokenId)

	tokenKeyLock := fmt.Sprintf("%s:Lock", tokenKey)
	boot.RedisCommonClient.SetLock(tokenKeyLock, 10000000)

	var accessToken string
	boot.RedisCommonClient.Get(tokenKey, &accessToken)
	if tool.IsEmpty(accessToken) {
		ret, err := refreshToken(region, tokenId)
		if err != nil {
			g.Log().Println(err, "------err1")
			return ""
		}

		// 缓存时间设置为access_token有效期-10秒
		boot.RedisCommonClient.SetEx(tokenKey, ret["access_token"], int(ret["expires_in"].(float64)-10))
		accessToken = ret["access_token"].(string)
	}

	// 释放锁
	_ = boot.RedisCommonClient.Del(tokenKeyLock)
	return accessToken
}

// 刷新AccessToken
func refreshToken(region RegionType, tokenId int) (map[string]interface{}, error) {
	cli := curl.NewHttpclient()
	var resToken string

	if tokenId == 0 {
		resToken = ResToken[region]
	} else {
		service := app.InitializeAuthService()
		recode, _ := service.SellerTokenRepo.GetOneData(tokenId)
		resToken = recode.RefreshToken
	}
	// 请求接口
	var postparama = make(map[string]string, 0)
	var res map[string]interface{}

	postparama["grant_type"] = "refresh_token"
	postparama["refresh_token"] = resToken
	postparama["client_id"] = ClientId
	postparama["client_secret"] = ClientSecret
	header := make(map[string]string, 0)
	header["User-Agent"] = "AdvertisingAPI PHP cli Library v1.2"
	dataflow := cli.IsJson(true).
		SetUrl(ApiTokenUrl).
		SetMethod("POST").
		GetDataFlow().
		SetWWWForm(postparama).
		SetHeader(header).
		BindJSON(&res)

	err := cli.Request(dataflow)

	if err != nil {
		g.Log().Error("Json Unmarshal err", err)
		return nil, err
	}

	// 状态码400以上表示请求错误
	if cli.Code >= 400 {
		g.Log().Error("http error", cli.Code)
		if res == nil {
			return nil, errors.New("the result is nil")
		} else {
			res["code"] = strconv.Itoa(cli.Code)
			if _, ok := res["message"]; ok {
				return nil, errors.New(res["message"].(string))
			} else {
				return nil, errors.New("the result  is nil1")
			}

		}

	}

	return res, nil
}
