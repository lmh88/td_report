// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// ReportSbConsumerDetail is the golang structure for table report_sb_consumer_detail.
type ReportSbConsumerDetail struct {
    Id         uint64      `orm:"id,primary"  json:"id"`         //                       
    CtxId      string      `orm:"ctx_id"      json:"ctxId"`      // 上下文链路追踪id      
    CreateTime *gtime.Time `orm:"create_time" json:"createTime"` // 创建时间              
    ReportName string      `orm:"report_name" json:"reportName"` //                       
    ProfileId  string      `orm:"profile_id"  json:"profileId"`  //                       
    Status     int         `orm:"status"      json:"status"`     // 0 默认 1 成功 2 失败  
    ReportDate string      `orm:"report_date" json:"reportDate"` //                       
    Error      string      `orm:"error"       json:"error"`      //                       
    Batch      string      `orm:"batch"       json:"batch"`      //                       
    CostTime   int         `orm:"cost_time"   json:"costTime"`   // 花费的时间秒          
    UpdateTime *gtime.Time `orm:"update_time" json:"updateTime"` // 修改时间              
}