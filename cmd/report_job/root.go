package report_job

import (
	"os"
	"td_report/cmd/report_job/job"
	"td_report/cmd/report_job/old"
	"td_report/cmd/report_job/task"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A brief description of your application",
	Long:  `总的根命令.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")
	RootCmd.PersistentFlags().StringVar(&old.StartDate, "startdate", "", "startdate:report start date")
	RootCmd.PersistentFlags().StringVar(&old.EndDate, "enddate", "", "enddate:report end date")
	RootCmd.PersistentFlags().StringVar(&old.ReportName, "report_name", "", "report_name:report_name")
	RootCmd.PersistentFlags().StringVar(&old.ReportType, "report_type", "", "report_type")
	RootCmd.PersistentFlags().StringVar(&old.ProfileId, "profile_id", "", "profile_id")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	var cmdList = []*cobra.Command{
		// 错误处理任务
		task.ErrortaskCmd,
		// 记录消费详情任务
		task.ConsumerDetailaskCmd,
		// s3上传任务
		task.UploadS3taskCmd,
		// s3 上传错误处理任务
		task.S3errorDealtaskCmd,
		// 新客户订阅fead消费脚本
		task.FeadtaskCmd,
		// dsp 切割文件
		old.DspDivideCmd,
		// dsp 拉取上个月的数据，时间周期是一段时间
		old.DspLastMonthCmd,
		// 新客户重新处理
		old.SchduleCmd,
		// 生产者
		old.NewproductCmd,
		// sb的消费者
		old.NewsbConsumerCmd,
		// sd的消费者
		old.NewsdConsumerCmd,
		// sp的消费者
		old.NewspConsumerCmd,
		// sb 品牌
		old.SbmetricsCmd,
		// dsp消费
		old.NewdspConsumerCmd,
		//fead订阅
		job.FeadSubCmd,
		// fead消费
		job.FeadConsumerCmd,
	}

	RootCmd.AddCommand(cmdList...)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".myapp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName("config.toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
