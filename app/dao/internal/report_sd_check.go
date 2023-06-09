// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// ReportSdCheckDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type ReportSdCheckDao struct {
	gmvc.M                       // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB               // DB is the raw underlying database management object.
	Table   string               // Table is the table name of the DAO.
	Columns reportSdCheckColumns // Columns contains all the columns of Table that for convenient usage.
}

// ReportSdCheckColumns defines and stores column names for table report_sd_check.
type reportSdCheckColumns struct {
	Id             string //
	ReportName     string //
	ReportDate     string //
	ProfileId      string //
	Filename       string //
	FileChangedate string // 文件更新时间
	Status         string // 0 不存在 1 存在
	Createdate     string //
	Updatedate     string //
	Extrant        string // 额外区分报表类型的字段，针对sd
}

var (
	// ReportSdCheck is globally public accessible object for table report_sd_check operations.
	ReportSdCheck = ReportSdCheckDao{
		M:     g.DB("report").Model("report_sd_check").Safe(),
		DB:    g.DB("report"),
		Table: "report_sd_check",
		Columns: reportSdCheckColumns{
			Id:             "id",
			ReportName:     "report_name",
			ReportDate:     "report_date",
			ProfileId:      "profile_id",
			Filename:       "filename",
			FileChangedate: "file_changedate",
			Status:         "status",
			Createdate:     "createdate",
			Updatedate:     "updatedate",
			Extrant:        "extrant",
		},
	}
)
