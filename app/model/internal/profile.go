// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// Profile is the golang structure for table profile.
type Profile struct {
    Id                  int64       `orm:"id,primary"          json:"id"`                  // 自增ID                             
    Region              string      `orm:"region"              json:"region"`              // 地区(NA:美区;EU:欧区;FE:日韩区)    
    ProfileId           int64       `orm:"profileId,unique"    json:"profileId"`           // ProfileID                          
    CountryCode         string      `orm:"countryCode"         json:"countryCode"`         // 国家代码                           
    CurrencyCode        string      `orm:"currencyCode"        json:"currencyCode"`        // 货币代码                           
    Timezone            string      `orm:"timezone"            json:"timezone"`            // 时区                               
    MarketplaceStringId string      `orm:"marketplaceStringId" json:"marketplaceStringId"` //                                    
    EntityId            string      `orm:"entityId"            json:"entityId"`            // 实体ID                             
    Type                string      `orm:"type"                json:"type"`                // 类型(seller:经销商;agency:代理商)  
    EntityName          string      `orm:"entityName"          json:"entityName"`          // 实体名称                           
    IsAuto              int         `orm:"isAuto"              json:"isAuto"`              //                                    
    Status              int         `orm:"status"              json:"status"`              // 状态(0:正常;1:禁用)                
    LastUpdate          *gtime.Time `orm:"lastUpdate"          json:"lastUpdate"`          // 最后更新时间                       
    AuthStatus          int         `orm:"authStatus"          json:"authStatus"`          // 同步状态                           
}