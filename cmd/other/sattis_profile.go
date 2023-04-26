package other

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var statisProfileCmd = &cobra.Command{
	Use:   "statis_profile",
	Short: "statis_profile",
	Long:  `statis_profile`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("statis_profile", false)
		logger.Logger.Info("statis_profile called")
		fmt.Println("test")
	},
}

func init() {
	RootCmd.AddCommand(statisProfileCmd)
}

func sp() {
	repo := repo.NewReportSpCheckRepository()
	data, err := repo.GetDataBycondition("2022-03-23", 1)
	if err != nil {
		fmt.Println(err, "=====")
	} else {

		if len(data) > 0 {
			total := 0
			checkNum := 0
			current := time.Now().Unix()
			chaMap := make(map[string]*model.ReportSpCheck)
			for _, item := range data {
				if item.FileChangedate != "" {
					total++
					if date, err := time.Parse(vars.TIMEFORMAT, item.FileChangedate); err != nil {
						fmt.Println(err.Error())
					} else {
						cha := current - date.Unix()
						if cha > 86400 {
							checkNum++
							key := fmt.Sprintf("%s_%s", item.ProfileId, item.ReportDate)
							chaMap[key] = item
						}
					}
				}
			}

			fmt.Println("total:", total)
			fmt.Println("cha:", len(chaMap))
			for _, it := range chaMap {
				fmt.Println(it.ProfileId, "===", it.ReportDate, "==", it.FileChangedate, "===", it.ReportName)
			}
		} else {
			fmt.Println("no data")
		}
	}
}

func sd() {
	repo := repo.NewReportSdCheckRepository()
	data, err := repo.GetDataBycondition("2022-03-23", 1)
	if err != nil {
		fmt.Println(err, "=====")
	} else {

		if len(data) > 0 {
			total := 0
			checkNum := 0
			current := time.Now().Unix()
			chaMap := make(map[string]*model.ReportSdCheck)
			for _, item := range data {
				if item.FileChangedate != "" {
					total++
					if date, err := time.Parse(vars.TIMEFORMAT, item.FileChangedate); err != nil {
						fmt.Println(err.Error())
					} else {
						cha := current - date.Unix()
						if cha > 86400 {
							checkNum++
							key := fmt.Sprintf("%s_%s", item.ProfileId, item.ReportDate)
							chaMap[key] = item
						}
					}
				}
			}

			fmt.Println("total:", total)
			fmt.Println("cha:", len(chaMap))
			for _, it := range chaMap {
				fmt.Println(it.ProfileId, "===", it.ReportDate, "==", it.FileChangedate, "===", it.ReportName, "===", it.Extrant)
			}
		} else {
			fmt.Println("no data")
		}
	}
}

func sb() {
	repo := repo.NewReportSbCheckRepository()
	data, err := repo.GetDataBycondition("2022-03-23", 1)
	if err != nil {
		fmt.Println(err, "=====")
	} else {

		if len(data) > 0 {
			total := 0
			checkNum := 0
			current := time.Now().Unix()
			chaMap := make(map[string]*model.ReportSbCheck)
			for _, item := range data {
				if item.FileChangedate != "" {
					total++
					if date, err := time.Parse(vars.TIMEFORMAT, item.FileChangedate); err != nil {
						fmt.Println(err.Error())
					} else {
						cha := current - date.Unix()
						if cha > 86400 {
							checkNum++
							key := fmt.Sprintf("%s_%s", item.ProfileId, item.ReportDate)
							chaMap[key] = item
						}
					}
				}
			}

			fmt.Println("total:", total)
			fmt.Println("cha:", len(chaMap))
			for _, it := range chaMap {
				fmt.Println(it.ProfileId, "===", it.ReportDate, "==", it.FileChangedate, "===", it.ReportName)
			}
		} else {
			fmt.Println("no data")
		}
	}
}

func dsp() {
	repo := repo.NewReportDspCheckRepository()
	data, err := repo.GetDataBycondition("2022-03-23", 1)
	if err != nil {
		fmt.Println(err, "=====")
	} else {

		if len(data) > 0 {
			total := 0
			checkNum := 0
			ioc := tool.GetIoc()
			current := time.Now().In(ioc).Unix()
			chaMap := make(map[string]*model.ReportDspCheck)
			for _, item := range data {
				if item.FileChangedate != "" {
					total++
					if date, err := time.Parse(vars.TIMEFORMAT, item.FileChangedate); err != nil {
						fmt.Println(err.Error())
					} else {
						cha := current - date.Unix()
						if cha > 86400 {
							checkNum++
							key := fmt.Sprintf("%s_%s", item.ProfileId, item.ReportDate)
							chaMap[key] = item
						}
					}
				}
			}

			fmt.Println("total:", total)
			fmt.Println("cha:", len(chaMap))
			for _, it := range chaMap {
				fmt.Println(it.ProfileId, "===", it.ReportDate, "==", it.FileChangedate, "===", it.ReportName)
			}
		} else {
			fmt.Println("no data")
		}
	}
}

func statisfunc() {
	dsp()
}
