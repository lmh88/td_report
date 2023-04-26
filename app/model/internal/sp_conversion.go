// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// SpConversion is the golang structure for table sp_conversion.
type SpConversion struct {
    Id                               uint        `orm:"id,primary"                            json:"id"`                                     //                         
    AdvertiserId                     string      `orm:"advertiser_id"                         json:"advertiser_id"`                          //                         
    MarketplaceId                    string      `orm:"marketplace_id"                        json:"marketplace_id"`                         //                         
    DatasetId                        string      `orm:"dataset_id"                            json:"dataset_id"`                             //                         
    IdempotencyId                    string      `orm:"idempotency_id,unique"                 json:"idempotency_id"`                         //                         
    AttributedSales1DSameSku         float64     `orm:"attributed_sales_1d_same_sku"          json:"attributed_sales_1d_same_sku"`          //
    AttributedConversions1D          int         `orm:"attributed_conversions_1d"             json:"attributed_conversions_1d"`             //
    AttributedSales1D                float64     `orm:"attributed_sales_1d"                   json:"attributed_sales_1d"`                   //
    AttributedConversions7D          int         `orm:"attributed_conversions_7d"             json:"attributed_conversions_7d"`             //
    AttributedConversions14DSameSku  int         `orm:"attributed_conversions_14d_same_sku"   json:"attributed_conversions_14d_same_sku"`   //
    AttributedConversions14D         int         `orm:"attributed_conversions_14d"            json:"attributed_conversions_14d"`            //
    TimeWindowStart                  *gtime.Time `orm:"time_window_start"                     json:"time_window_start"`                      // 特殊时间，带时区的时间  
    AttributedSales7D                float64     `orm:"attributed_sales_7d"                   json:"attributed_sales_7d"`                   //
    AttributedConversions30D         int         `orm:"attributed_conversions_30d"            json:"attributed_conversions_30d"`            //
    AttributedUnitsOrdered14DSameSku int         `orm:"attributed_units_ordered_14d_same_sku" json:"attributed_units_ordered_14d_same_sku"` //
    AttributedUnitsOrdered30D        int         `orm:"attributed_units_ordered_30d"          json:"attributed_units_ordered_30d"`          //
    AttributedSales7DSameSku         float64     `orm:"attributed_sales_7d_same_sku"          json:"attributed_sales_7d_same_sku"`          //
    AttributedUnitsOrdered14D        float64     `orm:"attributed_units_ordered_14d"          json:"attributed_units_ordered_14d"`          //
    AttributedUnitsOrdered7D         int         `orm:"attributed_units_ordered_7d"           json:"attributed_units_ordered_7d"`           //
    AttributedUnitsOrdered7DSameSku  int         `orm:"attributed_units_ordered_7d_same_sku"  json:"attributed_units_ordered_7d_same_sku"`  //
    AttributedConversions30DSameSku  int         `orm:"attributed_conversions_30d_same_sku"   json:"attributed_conversions_30d_same_sku"`   //
    AdGroupId                        string      `orm:"ad_group_id"                           json:"ad_group_id"`                            //                         
    Placement                        string      `orm:"placement"                             json:"placement"`                              //                         
    AttributedUnitsOrdered1D         int         `orm:"attributed_units_ordered_1d"           json:"attributed_units_ordered_1d"`           //
    AttributedSales30D               float64     `orm:"attributed_sales_30d"                  json:"attributed_sales_30d"`                  //
    AttributedUnitsOrdered30DSameSku int         `orm:"attributed_units_ordered_30d_same_sku" json:"attributed_units_ordered_30d_same_sku"` //
    AttributedUnitsOrdered1DSameSku  int         `orm:"attributed_units_ordered_1d_same_sku"  json:"attributed_units_ordered_1d_same_sku"`  //
    Currency                         string      `orm:"currency"                              json:"currency"`                               //                         
    AdId                             string      `orm:"ad_id"                                 json:"ad_id"`                                  //                         
    AttributedConversions1DSameSku   int         `orm:"attributed_conversions_1d_same_sku"    json:"attributed_conversions_1d_same_sku"`    //
    AttributedSales14D               float64     `orm:"attributed_sales_14d"                  json:"attributed_sales_14d"`                  //
    AttributedConversions7DSameSku   int         `orm:"attributed_conversions_7d_same_sku"    json:"attributed_conversions_7d_same_sku"`    //
    AttributedSales30DSameSku        float64     `orm:"attributed_sales_30d_same_sku"         json:"attributed_sales_30d_same_sku"`         //
    CampaignId                       string      `orm:"campaign_id"                           json:"campaign_id"`                            //                         
    KeywordId                        string      `orm:"keyword_id"                            json:"keyword_id"`                             //                         
    AttributedSales14DSameSku        float64     `orm:"attributed_sales_14d_same_sku"         json:"attributed_sales_14d_same_sku"`         //
    CreateDate                       *gtime.Time `orm:"create_date"                           json:"create_date"`                            //                         
}