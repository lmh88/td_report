// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// ReportBatchDetailDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type ReportBatchDetailDao struct {
	gmvc.M                           // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB                   // DB is the raw underlying database management object.
	Table   string                   // Table is the table name of the DAO.
	Columns reportBatchDetailColumns // Columns contains all the columns of Table that for convenient usage.
}

// ReportBatchDetailColumns defines and stores column names for table report_batch_detail.
type reportBatchDetailColumns struct {
	Id             string //
	Batch          string // 批次号，关联report_batch
	ReportNameList string // 同步的报表名称
	ReportType     string //
	StartDate      string //
	EndDate        string //
	Status         string // 1 执行中 2 完成 3 错误
	CreateDate     string //
	UpdateDate     string //
	Reason         string //
}

var (
	// ReportBatchDetail is globally public accessible object for table report_batch_detail operations.
	ReportBatchDetail = ReportBatchDetailDao{
		M:     g.DB("report").Model("report_batch_detail").Safe(),
		DB:    g.DB("report"),
		Table: "report_batch_detail",
		Columns: reportBatchDetailColumns{
			Id:             "id",
			Batch:          "batch",
			ReportNameList: "report_name_list",
			ReportType:     "report_type",
			StartDate:      "start_date",
			EndDate:        "end_date",
			Status:         "status",
			CreateDate:     "create_date",
			UpdateDate:     "update_date",
			Reason:         "reason",
		},
	}
)
