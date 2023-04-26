package other

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/container/gset"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"td_report/pkg/logger"
)

var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "temp",
	Long:  `temp`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("temp", true)
		logger.Logger.Info("temp called", true)
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListDspRegionProfile error",
			"err":  "error",
		})

		path1 := "D:\\data\\report\\new_dsp_back\\order\\20220101\\20220101_1898361487593447.csv"
		path2 := "D:\\data\\report\\new_dsp\\order\\20220201\\20220201_4109964880145.csv"
		readCsv(path1, path2)
	},
}

func init() {
	RootCmd.AddCommand(tempCmd)
}

func readCsv(path1, path2 string) {
	var line1 []string
	var line2 []string
	csvFile1, err := os.Open(path1)
	csvFile2, err := os.Open(path2)
	var s1 *gset.StrSet = gset.NewStrSet()
	var s2 *gset.StrSet = gset.NewStrSet()
	if err != nil {
		fmt.Println(err)
		return
	}

	reader1 := csv.NewReader(bufio.NewReader(csvFile1))
	reader2 := csv.NewReader(bufio.NewReader(csvFile2))
	reader2.Read()

	for {

		line1, err = reader1.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			s1.Add(line1...)
			break
		}
	}

	for {

		line2, err = reader2.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			s2.Add(line2...)
			break
		}
	}

	df1 := s2.Diff(s1)
	fmt.Println(df1.Size())
	df2 := s1.Diff(s2)
	fmt.Println(df2.Size())
}
