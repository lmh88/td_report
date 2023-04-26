package tools

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"td_report/app/repo"
	"td_report/vars"
	"time"
)

var statisDetailCmd = &cobra.Command{
	Use:   "statis_detail",
	Short: "statis_detail",
	Long:  `statis_detail`,
	Run: func(cmd *cobra.Command, args []string) {
		if StartDate == "" {
			StartDate = time.Now().Format(vars.TimeLayout)
		}

		statisDetailFunc(StartDate)
	},
}

func init() {
	RootCmd.AddCommand(statisDetailCmd)
}

type Mydata struct {
	profileId string
	nickname  string
}

func statisDetailFunc(startdate string) {
	reportTypeList := []string{vars.SP, vars.SB, vars.SD}
	data := make(map[string][]Mydata, 0)
	profileRepo := repo.NewSellerProfileRepository()
	profileMap, _ := profileRepo.GetAll()
	for _, reportType := range reportTypeList {
		path := vars.MypathMap[reportType]
		reportName := "campaigns"
		tempPath := path + "/" + reportName + "/" + startdate
		files, _ := ioutil.ReadDir(tempPath)
		key := fmt.Sprintf("%s_%s_%s", startdate, reportType, reportName)
		num := len(files)
		if num > 0 {
			MydataList := make([]Mydata, 0)
			for _, item := range files {
				if item.IsDir() {
					continue
				}
				arr := strings.Split(item.Name(), "_")
				if len(arr) > 1 {
					arr1 := strings.Split(arr[1], ".")
					profileIdStr := arr1[0]
					temp := Mydata{
						profileId: profileIdStr,
						nickname:  profileMap[profileIdStr],
					}
					MydataList = append(MydataList, temp)
				}
			}

			data[key] = MydataList
		}
	}

	file, err := os.Create(fmt.Sprintf("%s.csv", startdate))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	write := csv.NewWriter(file)
	write.Comma = ','
	for key, v := range data {
		for _, item := range v {
			myata := make([]string, 0)
			myata = append(myata, key)
			myata = append(myata, item.nickname)
			myata = append(myata, item.profileId)
			write.Write(myata)
		}
	}

	write.Flush()

}
