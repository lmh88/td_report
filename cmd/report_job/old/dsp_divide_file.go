package old

import (
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"os"
	"strings"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/service/report"
	"td_report/common/reporttool"
	"td_report/common/tool"
	"td_report/pkg/common"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var DspDivideCmd = &cobra.Command{
	Use:   "dsp_divide_file",
	Short: "dsp_divide_file",
	Long:  `dsp类型的将一个时间区间的文件切割成每天的文件`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("dsp_divide_file", false)
		logger.Logger.Info("dsp_divide_file called", time.Now().Format(vars.TIMEFORMAT))
		dspDivideFunc()
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 将一个时间周期内的文件切割成单一每天的文件，并且存入状态库
func dspDivideFunc() {
	sleepTimeNum := 1
	taskService := app.InitializeReportTaskService()
	commonpath := g.Cfg().GetString("common.datapath")
	csvService := report.NewDspCsvService()
	for {

		batch, err := reporttool.GetDivideFile(vars.DSP)
		if err != nil {
			logger.Logger.Info(map[string]interface{}{
				"flag": err.Error(),
			})

			if sleepTimeNum > 10 {
				sleepTime := time.Duration(rand.Int31n(1000))
				time.Sleep(sleepTime * time.Second)
				sleepTimeNum = 1
			} else {
				sleepTimeNum = sleepTimeNum + 1
				time.Sleep(10 * time.Second)
			}

			continue
		}

		strArr := strings.Split(batch, "_")
		// 目前不支持audience, audience数据内部是根据时间区间统计来的，没法按照每天的时间来拆分
		if strArr[0] == "audience" {
			continue
		}
		reportnameLen := len(strArr[0])
		fileName := batch[reportnameLen+1:]
		path := fmt.Sprintf("%s/%s/%s/%s", commonpath+vars.DspPathTemp, strArr[0], strArr[1], fileName)
		data := csvService.LoadCsvCfg(path, 1)
		if data != nil {
			daylists, err := tool.GetDaysDsp(strArr[1], strArr[2], vars.TimeLayout)
			if err != nil {
				logger.Error("dsp get date error ", err)
				continue
			}

			csvdata := make(map[string][][]string)
			Fileds := repo.GetDspFileFileds(strArr[0])
			for _, date := range daylists {
				startdate, _ := time.Parse(vars.TIMEDSPFILE, date)
				start := startdate.Format(vars.TimeLayout)
				dir := fmt.Sprintf("%s/%s/%s", commonpath+vars.DspPath, strArr[0], start)
				common.CreateDir(dir)
				filename := fmt.Sprintf("%s_%s", start, strArr[3])
				filepath := fmt.Sprintf("%s/%s", dir, filename)
				nfs, _ := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
				nfs.Seek(0, io.SeekEnd)
				w := csv.NewWriter(nfs)
				csvdata[date] = append(csvdata[date], Fileds)
				for _, item := range data.Records {
					if item.Record["date"] == date {
						csvdataList := make([]string, 0)
						for _, key := range Fileds {
							temp := item.Record[key]
							csvdataList = append(csvdataList, temp)
						}

						csvdata[date] = append(csvdata[date], csvdataList)

					}
				}

				w.WriteAll(csvdata[date])
				w.Flush()
				nfs.Close()
				taskService.Add(dir, strArr[0], vars.DSP, start)
			}
		}
	}
}
