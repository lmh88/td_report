// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// ProfileDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type ProfileDao struct {
	gmvc.M                                      // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB                              // DB is the raw underlying database management object.
	Table   string                              // Table is the table name of the DAO.
	Columns profileColumns // Columns contains all the columns of Table that for convenient usage.
}

// ProfileColumns defines and stores column names for table profile.
type profileColumns struct {
	Id                   string // 自增ID                             
    Region               string // 地区(NA:美区;EU:欧区;FE:日韩区)    
    ProfileId            string // ProfileID                          
    CountryCode          string // 国家代码                           
    CurrencyCode         string // 货币代码                           
    Timezone             string // 时区                               
    MarketplaceStringId  string //                                    
    EntityId             string // 实体ID                             
    Type                 string // 类型(seller:经销商;agency:代理商)  
    EntityName           string // 实体名称                           
    IsAuto               string //                                    
    Status               string // 状态(0:正常;1:禁用)                
    LastUpdate           string // 最后更新时间                       
    AuthStatus           string // 同步状态
}

var (
	// Profile is globally public accessible object for table profile operations.
	Profile = ProfileDao{
		M:     g.DB("default").Model("profile").Safe(),
		DB:    g.DB("default"),
		Table: "profile",
		Columns: profileColumns{
			Id:                  "id",                   
            Region:              "region",               
            ProfileId:           "profileId",            
            CountryCode:         "countryCode",          
            CurrencyCode:        "currencyCode",         
            Timezone:            "timezone",             
            MarketplaceStringId: "marketplaceStringId",  
            EntityId:            "entityId",             
            Type:                "type",                 
            EntityName:          "entityName",           
            IsAuto:              "isAuto",               
            Status:              "status",               
            LastUpdate:          "lastUpdate",           
            AuthStatus:          "authStatus",
		},
	}
)