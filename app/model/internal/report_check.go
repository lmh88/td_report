// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// ReportCheck is the golang structure for table report_check.
type ReportCheck struct {
    Id             string      `orm:"id,primary"      json:"id"`             //                                 
    ReportType     string      `orm:"report_type"     json:"reportType"`     //                                 
    ReportName     string      `orm:"report_name"     json:"reportName"`     //                                 
    ReportDate     *gtime.Time `orm:"report_date"     json:"reportDate"`     //                                 
    ProfileId      string      `orm:"profile_id"      json:"profileId"`      //                                 
    Filename       string      `orm:"filename"        json:"filename"`       //                                 
    FileChangedate string      `orm:"file_changedate" json:"fileChangedate"` // 文件更新时间                    
    Status         int         `orm:"status"          json:"status"`         // 0 不存在 1 存在                 
    Createdate     *gtime.Time `orm:"createdate"      json:"createdate"`     //                                 
    Updatedate     *gtime.Time `orm:"updatedate"      json:"updatedate"`     //                                 
    Extrant        string      `orm:"extrant"         json:"extrant"`        // 额外区分报表类型的字段，针对sd  
}