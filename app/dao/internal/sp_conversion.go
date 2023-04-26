// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// SpConversionDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type SpConversionDao struct {
	gmvc.M                                      // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB                              // DB is the raw underlying database management object.
	Table   string                              // Table is the table name of the DAO.
	Columns spConversionColumns // Columns contains all the columns of Table that for convenient usage.
}

// SpConversionColumns defines and stores column names for table sp_conversion.
type spConversionColumns struct {
	Id                                string //                         
    AdvertiserId                      string //                         
    MarketplaceId                     string //                         
    DatasetId                         string //                         
    IdempotencyId                     string //                         
    AttributedSales1DSameSku          string //                         
    AttributedConversions1D           string //                         
    AttributedSales1D                 string //                         
    AttributedConversions7D           string //                         
    AttributedConversions14DSameSku   string //                         
    AttributedConversions14D          string //                         
    TimeWindowStart                   string // 特殊时间，带时区的时间  
    AttributedSales7D                 string //                         
    AttributedConversions30D          string //                         
    AttributedUnitsOrdered14DSameSku  string //                         
    AttributedUnitsOrdered30D         string //                         
    AttributedSales7DSameSku          string //                         
    AttributedUnitsOrdered14D         string //                         
    AttributedUnitsOrdered7D          string //                         
    AttributedUnitsOrdered7DSameSku   string //                         
    AttributedConversions30DSameSku   string //                         
    AdGroupId                         string //                         
    Placement                         string //                         
    AttributedUnitsOrdered1D          string //                         
    AttributedSales30D                string //                         
    AttributedUnitsOrdered30DSameSku  string //                         
    AttributedUnitsOrdered1DSameSku   string //                         
    Currency                          string //                         
    AdId                              string //                         
    AttributedConversions1DSameSku    string //                         
    AttributedSales14D                string //                         
    AttributedConversions7DSameSku    string //                         
    AttributedSales30DSameSku         string //                         
    CampaignId                        string //                         
    KeywordId                         string //                         
    AttributedSales14DSameSku         string //                         
    CreateDate                        string //
}

var (
	// SpConversion is globally public accessible object for table sp_conversion operations.
	SpConversion = SpConversionDao{
		M:     g.DB("xray_report").Model("sp_conversion").Safe(),
		DB:    g.DB("xray_report"),
		Table: "sp_conversion",
		Columns: spConversionColumns{
			Id:                               "id",                                     
            AdvertiserId:                     "advertiser_id",                          
            MarketplaceId:                    "marketplace_id",                         
            DatasetId:                        "dataset_id",                             
            IdempotencyId:                    "idempotency_id",                         
            AttributedSales1DSameSku:         "attributed_sales_1d_same_sku",           
            AttributedConversions1D:          "attributed_conversions_1d",              
            AttributedSales1D:                "attributed_sales_1d",                    
            AttributedConversions7D:          "attributed_conversions_7d",              
            AttributedConversions14DSameSku:  "attributed_conversions_14d_same_sku",    
            AttributedConversions14D:         "attributed_conversions_14d",             
            TimeWindowStart:                  "time_window_start",                      
            AttributedSales7D:                "attributed_sales_7d",                    
            AttributedConversions30D:         "attributed_conversions_30d",             
            AttributedUnitsOrdered14DSameSku: "attributed_units_ordered_14d_same_sku",  
            AttributedUnitsOrdered30D:        "attributed_units_ordered_30d",           
            AttributedSales7DSameSku:         "attributed_sales_7d_same_sku",           
            AttributedUnitsOrdered14D:        "attributed_units_ordered_14d",           
            AttributedUnitsOrdered7D:         "attributed_units_ordered_7d",            
            AttributedUnitsOrdered7DSameSku:  "attributed_units_ordered_7d_same_sku",   
            AttributedConversions30DSameSku:  "attributed_conversions_30d_same_sku",    
            AdGroupId:                        "ad_group_id",                            
            Placement:                        "placement",                              
            AttributedUnitsOrdered1D:         "attributed_units_ordered_1d",            
            AttributedSales30D:               "attributed_sales_30d",                   
            AttributedUnitsOrdered30DSameSku: "attributed_units_ordered_30d_same_sku",  
            AttributedUnitsOrdered1DSameSku:  "attributed_units_ordered_1d_same_sku",   
            Currency:                         "currency",                               
            AdId:                             "ad_id",                                  
            AttributedConversions1DSameSku:   "attributed_conversions_1d_same_sku",     
            AttributedSales14D:               "attributed_sales_14d",                   
            AttributedConversions7DSameSku:   "attributed_conversions_7d_same_sku",     
            AttributedSales30DSameSku:        "attributed_sales_30d_same_sku",          
            CampaignId:                       "campaign_id",                            
            KeywordId:                        "keyword_id",                             
            AttributedSales14DSameSku:        "attributed_sales_14d_same_sku",          
            CreateDate:                       "create_date",
		},
	}
)