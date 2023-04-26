package tools

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"net/http"
	"runtime"
	"strings"
	"td_report/boot"
	"td_report/common/sendmsg"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"time"
)

var MonitorMqCmd = &cobra.Command{
	Use:   "monitor_mq",
	Short: "monitor_mq",
	Long:  `monitor_mq`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("monitor_mq", false)
		logger.Logger.Info("monitor_mq called")
		MonitorMqFunc()
	},
}

func init() {
	RootCmd.AddCommand(MonitorMqCmd)
}

// MonitorMqFunc 监控mq
func MonitorMqFunc() {
	var (
		username, password, urltemp, urlconst, url string
		checkMq                                    = true
		historyFlag                                = true
		monitorPort, host                          string //监控请求的端口
		path = "api/health/checks/local-alarms"
	)
	urlconst = "127.0.0.1"
	authstr := g.Cfg().GetString("rabbitmq.address")
	authResult := strings.Split(authstr, "@")
	sendObj:=sendmsg.New(sendmsg.Wechat)
	defer func() {
		if r := recover(); r != nil {
			sendObj.SendMsg("rabbitmq not connect !!!")
			restartSupervisor()
		}
	}()

	if len(authResult) > 1 {
		urlop := strings.Split(authResult[1], ":")
		urltemp = urlop[0]
		monitorPort = g.Cfg().GetString("rabbitmq.monitor_port")
		if strings.Contains(urltemp, urlconst) {
			if monitorPort == "" {
				host = urlconst
			} else {
				host = fmt.Sprintf("%s:%s", urlconst, monitorPort)
			}

			url = fmt.Sprintf("http://%s/%s", host, path)
		} else {

			if monitorPort == "" {
				host = urltemp
			} else {
				host = fmt.Sprintf("%s:%s", urltemp, monitorPort)
			}

			// 只是目前适合当前的url,不需要对应的端口
			if strings.Contains(authstr, "amqps") {
				url = fmt.Sprintf("https://%s/%s", host, path)
			} else {
				url = fmt.Sprintf("http://%s/%s", host,path)
			}
		}

		result := strings.Split(authResult[0], "//")
		if len(result) > 1 {
			so := strings.Split(result[1], ":")
			username = so[0]
			password = so[1]
		}
	}

	client := g.Client()
	client = client.SetBasicAuth(username, password)
	n := 5
	// 如果重试n次还是错误
	for i := 0; i < n; i++ {
		repos, err := client.Get(url)
		if err != nil {
			g.Log().Error(err.Error())
			panic(err)
			checkMq = false
			historyFlag = false
		} else {

			// 如果请求不是200 ，则当前健康监测
			if repos.StatusCode != http.StatusOK {
				checkMq = false
				historyFlag = false

			} else {

				if _, err := boot.GetRabbitmqClient(); err != nil {
					checkMq = false
					historyFlag = false
				} else {
					checkMq = true
					break
				}
			}
		}

		time.Sleep(2 * time.Second)
	}

	if checkMq == false {
		// 1、发送消息到企业微信
		// 2、重启对应的supervisor进程
		sendObj.SendMsg("rabbitmq not connect !!!")
		restartSupervisor()
	} else {
		//// 如果曾今mq链接补上，但是后期链接上了，还是需要重启服务
		if historyFlag == false {
			restartSupervisor()
			fmt.Println("history check error sendmsg")
		}
	}
}

//重启supervisor 对应的进程
func restartSupervisor() {
	//如果是windows，直接跳过
	if runtime.GOOS == "windows" {
		return
	}

	processList := []string{
		"dsp_consumer:consumer_00",
		"dsp_consumer:consumer_01",
		"dsp_consumer:consumer_02",
		"dsp_consumer:consumer_03",
		"dsp_consumer:consumer_04",
		"dsp_consumer:consumer_05",
		"dsp_slow_consumer:consumer_00",
		"err_task:consumer_00",
		"s3_upload:consumer_00",
		"s3_upload:consumer_01",
	}

	cmd := "supervisorctl"
	path := g.Cfg().GetString("rabbitmq.execpath")
	for _, item := range processList {
		args:= make([]string, 0)
		args = append(args, fmt.Sprintf("%s %s", "restart", item))
		tool.RunCommand(path, cmd, args...)
	}
}
