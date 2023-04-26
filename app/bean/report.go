package bean

import "github.com/gogf/gf/os/gtime"

type ReportCommonData struct {
	ReportName    string   `json:"reportName"`
	ReportType    string   `json:"reportType"`
	StartDate     string   `json:"startDate"`
	EndDate       string   `json:"endDate"`
	ProfileIdList []string `json:"profile_list"`
}

// Tempdata 记录check后重试数据，extrant 主要针对sd标志, t200 t300
type Tempdata struct {
	ProfileId  string
	Extrant    string
	ErrorType  int //错误类型 1 文件为空 2 文件没有及时更新
	FileUpdate string
}
type Ckeckdata struct {
	ProfileId  string
	Extrant    string
	ErrorType  int //错误类型 1 文件为空 2 文件没有及时更新
	FileUpdate string
	ReportDate string
}

type CheckDetail struct {
	ReportType string      `orm:"report_type"   json:"reportType"` //
	ReportName string      `orm:"report_name"   json:"reportName"` //
	ReportDate string      `orm:"report_date"   json:"reportDate"` //
	ProfileId  string      `orm:"profile_id"    json:"profileId"`  //
	RetryType  int         `orm:"retry_type"    json:"retryType"`  // 1 文件为空 2 文件没更新 0 默认
	CreateDate *gtime.Time `orm:"create_date"   json:"createDate"` //
	UpdateDate *gtime.Time `orm:"update_date"   json:"updateDate"` //
	CheckTime  int         `orm:"check_time"    json:"checkTime"`  // 检测次数
	Extrant    string      `orm:"extrant"       json:"extrant"`    // 针对sd的t200 或者t300等类别
	FileUpdate string      `orm:"file_update"   json:"file_update"`
}

type Done struct {
	ReportType string
	Date       string
}

type ReportErr struct {
	ReportType  string
	ReportName  string
	ProfileId   string
	ReportDate  string
	ErrorType   int
	ErrorReason string
	KeyParam    string
	Extra       string
}

