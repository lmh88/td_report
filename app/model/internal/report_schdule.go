// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// ReportSchdule is the golang structure for table report_schdule.
type ReportSchdule struct {
    Id             int         `orm:"id,primary"      json:"id"`             //                                                                      
    Batch          string      `orm:"batch"           json:"batch"`          //                                                                      
    ReportType     string      `orm:"report_type"     json:"reportType"`     //                                                                      
    ReportnameList string      `orm:"reportname_list" json:"reportnameList"` //                                                                      
    ReportDate     string      `orm:"report_date"     json:"reportDate"`     //                                                                      
    CreateDate     *gtime.Time `orm:"create_date"     json:"createDate"`     //                                                                      
    ProfileNum     int         `orm:"profile_num"     json:"profileNum"`     //                                                                      
    SchduleType    int         `orm:"schdule_type"    json:"schduleType"`    // 调度类型 0:今天和昨天 1:14天内 2: 半月调度3:整个月调度4检测文件调度
    EndTime        *gtime.Time `orm:"end_time"  json:"endTime"`// 结束时间
    StartTime        *gtime.Time `orm:"start_time"  json:"startTime"`// 开始时间
}