package old

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/common/reportsystem"
	"td_report/common/reporttool"
	"td_report/pkg/dsp/dsp_all_steps"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// DspLastMonthCmd 针对dsp在汇率下来后拉取上个月的全部数据
var DspLastMonthCmd = &cobra.Command{
	Use:   "dsp_last_month",
	Short: "dsp_last_month",
	Long:  `dsp_last_month`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("dsp_last_month", false)
		logger.Logger.Info("dsp_last_month called")
		StartDate = "20220523"
		EndDate = "20220525"
		if StartDate == "" || EndDate == "" {
			logger.Logger.Error(map[string]interface{}{
				"flag":      "dsp last month params error",
				"startdate": StartDate,
				"enddate":   EndDate,
			})

			return
		}

		dspLastMonthFunc(StartDate, EndDate)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func dspLastMonthFunc(start, end string) {
	profileRepository := repo.NewProfileRepository()
	profileList := []string{"1430615438938675", "2997968305172217"}

	rps, err := profileRepository.ListDspRegionProfile(profileList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListDspProfile error",
			"err":  err.Error(),
		})
		return
	}

	n := 100 // 数据切割成按照100一份
	groupData := profileRepository.ArrayInGroupsOf(rps, int64(n))
	dspReportNameList := vars.ReportList[vars.DSP]
	wg := reportsystem.NewPool(2)
	dayListMap, err := reporttool.GetDspPeriod(start, end, vars.DSP)
	for _, reportname := range dspReportNameList {
		// audience 不支持此模式，里面的数据没办法拆分
		if reportname == "audience" {
			continue
		}
		for _, date := range dayListMap[reportname] {
			for _, v := range groupData {
				for _, item := range v {
					wg.Add(1)
					go func(gitem *bean.DspRegionProfile, greportname string, gDate *bean.Mydate) {
						defer wg.Done()
						ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.DSP, greportname, gitem.ProfileId, time.Now().Unix()))
						dsp_all_steps.DspAllStepsPeriod(gitem.Region, gitem.ProfileId, greportname, gDate.StartDate, gDate.EndDate, ctx)
					}(item, reportname, date)
				}
			}
		}
	}

	wg.Wait()
	logger.Logger.Info("pull done")
}
