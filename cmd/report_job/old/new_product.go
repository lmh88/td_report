package old

import (
	"github.com/spf13/cobra"
	"runtime"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/server"
	"td_report/app/service/report"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// productday 生产的天数 2:今天和昨天 7:最近7天  14: 最近14天
var productday int

// 0 默认传参数按照天数来传输; 1 天数忽略使用的是全月，每月5号执行
var productType string

// 目前暂定快速拉取的天数是2 ，需要过滤，其他的不需要过滤
const (
	// 特定天数，如果是2天的情况就需要过滤黑名单
	spevialDay = 2
)

// NewproductCmd 新的生产者，只是产生日期，不做其他的处理，新模式下的消费者获取对应的日期异步消费
var NewproductCmd = &cobra.Command{
	Use:   "new_product",
	Short: "new_product",
	Long:  `new_product`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("new_product", false)
		logger.Logger.Info("new_product called", time.Now().Format(vars.TIMEFORMAT))
		if productday == 0 {
			logger.Logger.Info("生产的天数错误，不能为0")
			return
		}

		var currentTime time.Time
		if runtime.GOOS == "windows" { // windows 没有对应的时区文件，获取时间错误，默认获取本地时间
			currentTime = time.Now()
		} else {

			localzon,err:= time.LoadLocation("Local")
			if err != nil {
				logger.Logger.Info("获取本地时区错误")
				return
			}
			currentTime = time.Now().In(localzon)
		}

		currentHouer := currentTime.Hour()
		//如果是北京时间晚上21点到第二天的8点，则不执行快速的生产者，跳过，留给全量跑
		if productday == vars.ProductDay2 {
			if (currentHouer >= 21 && currentHouer <= 24) || (currentHouer >= 0 && currentHouer <= 8) {
				logger.Logger.Info("快速通道当前不能生产数据跳过")
				return
			}
		}

		profileService := app.InitializeProfileService()
		schduleRepo := repo.NewReportSchduleRepository()
		retryRepo := repo.NewReportCheckRetryDetailRepository()
		addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
		productServer := server.NewProductServer(addQueueService)
		productServer.Product(productday, ReportType, spevialDay)
	},
}

func init() {
	NewproductCmd.PersistentFlags().IntVar(&productday, "day", 2, "生产的天数")
	NewproductCmd.PersistentFlags().StringVar(&productType, "product_type", "0", "生产是按照天数计算还是按照整个月计算:0 默认按照天数;1按照整个月")
}
