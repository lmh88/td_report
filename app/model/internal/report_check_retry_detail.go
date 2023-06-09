// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// ReportCheckRetryDetail is the golang structure for table report_check_retry_detail.
type ReportCheckRetryDetail struct {
    Id           string      `orm:"id,primary"    json:"id"`           //                                                               
    ReportType   string      `orm:"report_type"   json:"reportType"`   //                                                               
    ReportName   string      `orm:"report_name"   json:"reportName"`   //                                                               
    ReportDate   string      `orm:"report_date"   json:"reportDate"`   //                                                               
    ProfileId    string      `orm:"profile_id"    json:"profileId"`    //                                                               
    RetryType    int         `orm:"retry_type"    json:"retryType"`    // 1 文件为空 2 文件没更新 0 默认                                
    CreateDate   *gtime.Time `orm:"create_date"   json:"createDate"`   //                                                               
    UpdateDate   *gtime.Time `orm:"update_date"   json:"updateDate"`   //                                                               
    CheckTime    int         `orm:"check_time"    json:"checkTime"`    // 检测次数                                                      
    Extrant      string      `orm:"extrant"       json:"extrant"`      // 针对sd的t200 或者t300等类别                                   
    ErrorsReason string      `orm:"errors_reason" json:"errorsReason"` // 错误情况，针对检测程序发现重拾次数太多后拉取部分记录错误情况  
}