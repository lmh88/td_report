package check_report_file

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	rediS "github.com/go-redis/redis"
	"io/ioutil"
	"os"
	"strings"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

const Prefix = "check_file"

func AddRedisSetVal(rds *rediS.Client, key string, data []string, expire int) {
	rds.Del(key)
	rds.SAdd(key, data)
	rds.Expire(key, time.Duration(expire)*time.Second)
}

func CountGzipLen(tarName string) int {
	if tarName == "" {
		return 1
	}
	tarFile, err := os.Open(tarName)
	defer tarFile.Close()
	if err != nil {
		logger.Logger.Error(tarName, "未找到, error:", err.Error())
		return 1
	}
	gr, err := gzip.NewReader(tarFile)
	if err != nil {
		logger.Logger.Error(tarName, "error:", err.Error())
		return 1
	}
	con, err := ioutil.ReadAll(gr)
	if err != nil {
		logger.Logger.Error(tarName, "error:", err.Error())
		return 1
	}
	tmp := make([]interface{}, 0)
	err = json.Unmarshal(con, &tmp)
	if err != nil {
		logger.Logger.Error(tarName, "json_error:", err.Error())
		return 1
	}
	return len(tmp)
}

func GetFileProfile(fileName string) string {
	strArr := strings.Split(fileName, ".")
	strArr2 := strings.Split(strArr[0], "_")
	return strArr2[1]
}

func CountFileLine(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer file.Close()
	fd := bufio.NewReader(file)
	count := 0
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	return count
}

func GetFilesByPath(path string) (files []string) {
	files = make([]string, 0)
	fArr, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Logger.Error("路径错误：", path, ",error:", err.Error())
		return
	}
	for _, file := range fArr {
		if file.IsDir() {
			continue
		} else {
			files = append(files, file.Name())
		}
	}
	//logger.Logger.Info(path, "获取文件", files)
	return
}

func GetLatelyThreeDay(date time.Time) []string {
	days := make([]string, 0)
	days = append(days, date.Format(vars.TimeLayout))
	one, _ := time.ParseDuration("-24h")
	days = append(days, date.Add(one).Format(vars.TimeLayout))
	one, _ = time.ParseDuration("-48h")
	days = append(days, date.Add(one).Format(vars.TimeLayout))
	//one, _ = time.ParseDuration("-72h")
	//days = append(days, date.Add(one).Format(vars.TimeLayout))
	return days
}

func GetTypePrefix(reportType string) string {
	return fmt.Sprintf("%s:%s:", Prefix, reportType)
}

func GetRdsKey(prefix string, date time.Time) string {
	return prefix + "merge_" + date.Format(vars.TimeLayout)
}

func GetRdsKeyOfdate(reportType string, date string) string {
	return fmt.Sprintf("%smerge_%s", GetTypePrefix(reportType), date)
}
