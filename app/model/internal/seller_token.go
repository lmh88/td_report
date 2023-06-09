// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// SellerToken is the golang structure for table seller_token.
type SellerToken struct {
    Id             int         `orm:"id,primary"     json:"id"`             //                                               
    UserId         int         `orm:"userId"         json:"userId"`         //                                               
    Region         string      `orm:"region"         json:"region"`         //                                               
    Country        string      `orm:"country"        json:"country"`        // 已勾选国家                                    
    SellerSN       string      `orm:"sellerSN"       json:"sellerSN"`       //                                               
    Nickname       string      `orm:"nickname"       json:"nickname"`       //

    ClientId       int         `orm:"clientId"       json:"clientId"`       // clientId

    RefreshToken   string      `orm:"refreshToken"   json:"refreshToken"`   //                                               
    AccessToken    string      `orm:"accessToken"    json:"accessToken"`    //                                               
    Expires        int         `orm:"expires"        json:"expires"`        //                                               
    SpRefreshToken string      `orm:"spRefreshToken" json:"spRefreshToken"` // SPAPI refresh token                           
    SpAccessToken  string      `orm:"spAccessToken"  json:"spAccessToken"`  // SPAPI access token                            
    SpExpires      int         `orm:"spExpires"      json:"spExpires"`      // SPAPI过期时间                                 
    Code           string      `orm:"code"           json:"code"`           //                                               
    Params         string      `orm:"params"         json:"params"`         //                                               
    Status         int         `orm:"status"         json:"status"`         // 1 - 启用，2 - 删除，3 - 暂存                  
    IsPPC          int         `orm:"isPPC"          json:"isPPC"`          // ppc授权状态，0：未授权 1：已授权              
    IsSP           int         `orm:"isSP"           json:"isSP"`           // SPAPI授权状态，0：未授权 1：已授权 2：已过期  
    IsMWS          int         `orm:"isMWS"          json:"isMWS"`          // mws授权状态，0：未授权 1：已授权              
    CreatedAt      *gtime.Time `orm:"createdAt"      json:"createdAt"`      //                                               
    LastUpdate     *gtime.Time `orm:"lastUpdate"     json:"lastUpdate"`     //                                               
}