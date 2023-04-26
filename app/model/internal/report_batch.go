// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// ReportBatch is the golang structure for table report_batch.
type ReportBatch struct {
    Id         uint        `orm:"id,primary"  json:"id"`         //                            
    Batch      string      `orm:"batch"       json:"batch"`      // 批次号                     
    Paramas    string      `orm:"paramas"     json:"paramas"`    // 参数                       
    Status     int         `orm:"status"      json:"status"`     // 状态 1 创建 2 成功 3 失败  
    CreateTime *gtime.Time `orm:"create_time" json:"createTime"` // 创建时间                   
    UpdateTime *gtime.Time `orm:"update_time" json:"updateTime"` // 修改时间                   
    IsCheck    int         `orm:"is_check"    json:"isCheck"`    // 默认0 已处理1              
}