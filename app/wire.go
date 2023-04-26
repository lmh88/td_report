//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package app

import (
	"github.com/google/wire"
	"td_report/app/repo"
	"td_report/app/service/amazon"
	"td_report/app/service/common"
	"td_report/app/service/report"
	"td_report/app/service/tool"
)

//初始化权限server
func InitializeAuthService() *common.AuthService {
	panic(wire.Build(common.AuthServiceSet,
		repo.NewSellerProfileRepository,
		repo.NewSellerTokenRepository,
		repo.NewProfileRepository,
	))
}

// 初始化报表任务服务
func InitializeReportTaskService() *common.ReportTaskService {
	panic(wire.Build(common.ReportTaskServiceSet,
		repo.NewReportTaskRepository,
	))
}

//  获取profile
func InitializeProfileService() *amazon.ProfileService {
	panic(wire.Build(amazon.ProfileServiceSet,
		repo.NewProfileRepository,
		repo.NewSellerProfileRepository,
	))
}

//初始化统计工具服务
func InitializeStatisticsService() *tool.StatisticsService {
	panic(wire.Build(tool.StatisticsServiceSet,
		repo.NewSellerProfileRepository,
		repo.NewProfileRepository,
		repo.NewReportSdCheckRepository,
		repo.NewReportSpCheckRepository,
		repo.NewReportSbCheckRepository,
		repo.NewReportDspCheckRepository,
	))
}

//初始化队列报表调度服务
func InitializeAllReportService() *report.AllReportService {
	panic(wire.Build(report.AllReportServiceSet,
		repo.NewSellerProfileRepository,
		repo.NewProfileRepository,
		repo.NewReportTaskRepository,
	))
}

// web入口授权等服务
func InitializeToolReportService() *report.ToolReportService {
	panic(wire.Build(report.ToolReportServiceSet,
		repo.NewSellerProfileRepository,
		repo.NewProfileRepository,
		repo.NewReportTaskRepository,
		repo.NewReportBatchRepository,
		repo.NewReportBatchDetailRepository,
	))
}

//新客户授权
func InitializeCustomerReportService() *report.CustomerReportService {
	panic(wire.Build(report.CustomerReportServiceSet,
		repo.NewSellerProfileRepository,
		repo.NewProfileRepository,
		repo.NewReportTaskRepository,
		repo.NewReportBatchRepository,
		repo.NewReportBatchDetailRepository,
	))
}

func InitializeFeadService() *report.FeadService {
	panic(wire.Build(report.FeadServiceSet,
		repo.NewSellerProfileRepository,
		repo.NewSpConversionRepository,
		repo.NewSpTrafficRepository,
		repo.NewFeadRepository,
	))
}
