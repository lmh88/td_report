package tools

import (
	"fmt"
	rediS "github.com/go-redis/redis"
	"github.com/spf13/cobra"
	"sync"
	"td_report/app/service/check_report_file"
	"td_report/boot"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var checkDate string

const (
	TenMinute = 600
	OneHour   = 3600
	OneDay    = 3600 * 24
)

var checkReportFile = &cobra.Command{
	Use:   "check_report_file",
	Short: "检查报表文件为空的profile",
	Long:  `check_report_file`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("check_report_file", false)
		logger.Logger.Info("check_report_file called")
	},
	Run: func(cmd *cobra.Command, args []string) {

		var (
			date time.Time
			err  error
		)

		if checkDate != "" {
			if date, err = time.Parse(vars.TimeLayout, checkDate); err != nil {
				logger.Logger.Error(map[string]interface{}{
					"err":  err,
					"desc": "时间格式错误",
				})
				fmt.Println("时间格式错误")
				return
			}
		} else {
			date = time.Now()
		}

		if ReportType == "" {
			fmt.Println("report_type必填")
			return
		}
		checkFile(date, ReportType)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("check_report_file run over")
	},
}

func init() {
	RootCmd.AddCommand(checkReportFile)
	RootCmd.PersistentFlags().StringVar(&checkDate, "check_date", "", "检查的日期，如：20220322")
}

func checkFile(date time.Time, reportType string) {

	prefix := check_report_file.GetTypePrefix(reportType)
	rds := boot.RedisCommonClient.GetClient()

	days := check_report_file.GetLatelyThreeDay(date)
	logger.Logger.Info(reportType, "获取日期", days)
	groupKeys := make([]string, 0)
	var wg sync.WaitGroup
	for _, day := range days {
		key := prefix + day
		groupKeys = append(groupKeys, key)
		wg.Add(1)
		go func(day, key, reportType string, rds *rediS.Client) {
			defer wg.Done()
			logger.Logger.Info(key, "开始执行")
			taskKeys := make([]string, 0)
			for _, taskType := range vars.ReportList[reportType] {
				taskKey := key + "_" + taskType
				taskKeys = append(taskKeys, taskKey)
				path := fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], taskType, day)
				files := check_report_file.GetFilesByPath(path)
				if len(files) == 0 {
					continue
				}
				profiles := make([]string, 0)
				profiles1 := make([]string, 0)
				profiles2 := make([]string, 0)
				for _, file := range files {
					if reportType == vars.DSP {
						if check_report_file.CountFileLine(path+"/"+file) == 1 {
							profiles1 = append(profiles1, check_report_file.GetFileProfile(file))
						} else {
							profiles2 = append(profiles2, check_report_file.GetFileProfile(file))
						}
					} else {
						if check_report_file.CountGzipLen(path+"/"+file) == 0 {
							profiles1 = append(profiles1, check_report_file.GetFileProfile(file))
						} else {
							profiles2 = append(profiles2, check_report_file.GetFileProfile(file))
						}
					}
				}
				//rds.SAdd(taskKey+"1", profiles1)
				check_report_file.AddRedisSetVal(rds, taskKey+"1", profiles1, TenMinute)
				//rds.SAdd(taskKey+"2", profiles2)
				check_report_file.AddRedisSetVal(rds, taskKey+"2", profiles2, TenMinute)
				profiles = rds.SDiff(taskKey+"1", taskKey+"2").Val()
				//rds.SAdd(taskKey, profiles)
				check_report_file.AddRedisSetVal(rds, taskKey, profiles, TenMinute)
			}
			res := rds.SInter(taskKeys...).Val()
			logger.Logger.Info(key, "结果", res, len(res))
			//rds.SAdd(key, res)
			check_report_file.AddRedisSetVal(rds, key, res, OneHour)
		}(day, key, reportType, rds)
	}
	wg.Wait()
	mergeProfile := rds.SInter(groupKeys...).Val()

	mergeKey := check_report_file.GetRdsKey(prefix, date)
	check_report_file.AddRedisSetVal(rds, mergeKey, mergeProfile, OneDay)
	logger.Logger.Info(mergeKey, "合并结束", mergeProfile, len(mergeProfile))
}
