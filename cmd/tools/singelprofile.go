package tools

import (
	"context"
	"fmt"
	"math/rand"
	"td_report/app"
	"td_report/common/reportsystem"
	"td_report/pkg/all_steps"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"

	"github.com/spf13/cobra"
)

// singelCmd represents the statisticsTool command
// 拉取单个profile
var singelprofileCmd = &cobra.Command{
	Use:   "singel_profile",
	Short: "singel_profile",
	Long:  `singel_profile`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("singel_profile", false)
		logger.Logger.Info("singel_profile called")
		ReportName = "campaigns"
		ReportType = "sb"
		StartDate = "20220902"
		ProfileId = "2697320314504359"
		if ReportType == "" || ProfileId == "" || ReportName == "" || StartDate == "" {
			logger.Logger.Info("singel paramas empty")
			return
		}

		singelProfileFunc(ReportType)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RootCmd.AddCommand(singelprofileCmd)
}

func singelProfileFunc(ReportType string) {
	profileService := app.InitializeProfileService()
	profileIdList := make([]string, 0)
	profileIdList = append(profileIdList, ProfileId)
	rps, err := profileService.GetPpcProfile(profileIdList)
	if err != nil || len(rps) == 0 {
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(rps) == 0 {
			fmt.Println("zero")
		}
		fmt.Println(" not found profile")
		return
	}

	rpt := rps[0]
	pool := reportsystem.NewPool(15)
	patch := time.Now().Unix()
	switch ReportType {
	case vars.SP:
		pool.Add(1)
		ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SP, ReportName, rpt.ProfileId, patch))
		go all_steps.SPAllSteps(rpt, ReportName, StartDate, pool, ctx)
		//todo 权重写入
	case vars.SD:
		for _, tactic := range []string{"T00020", "T00030"} {
			pool.Add(1)
			ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%s_%d", vars.SD, ReportName, rpt.ProfileId, tactic, patch))
			go all_steps.SDAllSteps(rpt, ReportName, StartDate, tactic, pool, ctx)
		}
	case vars.SB:
		pool.Add(1)
		ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SB, ReportName, rpt.ProfileId, patch))
		go all_steps.SbAllSteps(rpt, ReportName, StartDate, EndDate, pool, ctx)

	default:
		logger.Logger.Error(map[string]interface{}{
			"flag":       "reportType not found",
			"reportType": ReportType,
		})
		return
	}

	pool.Wait()
}
