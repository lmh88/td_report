package reporttool

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"td_report/app"
	"td_report/app/bean"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
)

func GetPath(reportType string) string {
	return vars.PathMap[reportType]
}

func GetDspPeriod(start, end, reportType string) (map[string][]*bean.Mydate, error) {
	result := make(map[string][]*bean.Mydate, 0)
	reportNameList := vars.ReportList[reportType]
	var period int
	for _, reportName := range reportNameList {
		reportKey := fmt.Sprintf("report.dsp_%s_length", reportName)
		period = g.Cfg().GetInt(reportKey)
		daylist, _ := tool.GetDaysPeriod(start, end, vars.TimeLayout, period)
		result[reportName] = daylist
	}

	return result, nil
}

// Getbatch 获取批次号
func Getbatch(reportType string) (string, error) {
	batchdata, err := boot.RedisCommonClient.GetClient().RPop(redis.WithBatch(reportType)).Result()
	if err != nil && err != redis.Nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "redis.Client.RPop batch error",
			"err":  err,
		})

		return "", err
	}

	if err == redis.Nil {
		logger.Logger.Info(map[string]interface{}{
			"flag": "task all batch done get patch",
			"err":  err,
		})

		return "", redis.Nil
	} else {

		return batchdata, nil
	}
}

func GetNewBatch(reportType string, QueuePrefix string) (string, error) {
	batchdata, err := boot.RedisCommonClient.GetClient().RPop(redis.WithNewBatch(reportType, QueuePrefix)).Result()
	if err != nil && err != redis.Nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "redis.Client.RPop batch error",
			"err":  err,
		})

		return "", err
	}

	if err == redis.Nil {
		logger.Logger.Info(map[string]interface{}{
			"flag": "task all batch done get patch",
			"err":  err,
		})

		return "", redis.Nil
	} else {

		return batchdata, nil
	}
}

func GetDivideFileHook(reportType string) (string, error) {
	return "order_20220201_20220216_1898361487593447.csv", nil
}

func GetDivideFile(reportType string) (string, error) {
	batchdata, err := boot.RedisCommonClient.GetClient().RPop(redis.WithDivide(reportType)).Result()
	if err != nil && err != redis.Nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "redis.Client.RPop batch error",
			"err":  err,
		})

		return "", err
	}

	if err == redis.Nil {
		logger.Logger.Info(map[string]interface{}{
			"flag": "divilede file all batch done get patch",
			"err":  err,
		})

		return "", redis.Nil
	} else {

		return batchdata, nil
	}
}

func GetPtsMap() (map[string]*bean.ProfileToken, error) {
	profileService := app.InitializeProfileService()
	profileIdList := make([]string, 0)
	ptsMap, err := profileService.GetPpcProfileMap(profileIdList)
	if err != nil {
		logger.Logger.Error(err, "get  profile error")
		return nil, err
	}
	return ptsMap, nil
}

func GetQueueLength(reportType string) (int64, error) {
	return boot.RedisCommonClient.GetClient().LLen(redis.WithBatch(reportType)).Result()
}

// CheckQueueNum false 可以添加 true 不能添加
func CheckQueueNum(reportType string, num int64) bool {
	length, err := boot.RedisCommonClient.GetClient().LLen(redis.WithBatch(reportType)).Result()
	if err != nil {
		logger.Logger.Error(err)
		return false
	}

	if length > num {
		logger.Logger.Info(reportType, "队列消费阻塞严重，暂时停止调度添加元素")
		return true
	}

	return false
}

func CheckQueueNewbatchNum(reportType string, num int64, QueuePrefix string) bool {
	length, err := boot.RedisCommonClient.GetClient().LLen(redis.WithNewBatch(reportType, QueuePrefix)).Result()
	if err != nil {
		logger.Logger.Error(err)
		return false
	}

	if length > num {
		logger.Logger.Info(reportType, "队列消费阻塞严重，暂时停止调度添加元素")
		return true
	}

	return false
}

func GetQueuePrefix(day int) string {
	var prefixQueue string
	// 根据不同天数添加队列前缀拉区分队列
	if day <= vars.ProductDay2 {
		prefixQueue = vars.FastQueue
	} else if day <= vars.ProductDay7 {
		prefixQueue = vars.MiddleQueue
	} else if day <= vars.ProductDay14 {
		prefixQueue = vars.SlowQueue
	} else {
		prefixQueue = vars.BackQueue
	}

	return prefixQueue
}

// CheckQueueName 检测输入是否正确
func CheckQueueName(queueType string) bool {
	for _, item := range vars.RedisQueueList {
		if item == queueType {
			return true
		}
	}

	return false
}
