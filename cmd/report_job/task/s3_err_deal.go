package task

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/bean"
	"td_report/boot"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/common/redis"
	"td_report/common/sendmsg/wechart"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"
)

var s3date string

// S3errorDealtaskCmd 处理s3 上传失败的情况
var S3errorDealtaskCmd = &cobra.Command{
	Use:   "s3err_deal",
	Short: "s3err_deal",
	Long:  `s3err_deal`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("s3err_deal", false)
		logger.Logger.Info("s3err_deal called", time.Now().Format(vars.TIMEFORMAT))
		s3errorFunc()
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func s3errorFunc() {
	s3key := redis.WithUploadFile("uploaderr")
	key := g.Cfg().GetString("wechat.key")
	send := wechart.NewSendMsg(key, g.Cfg().GetBool("wechat.open"))
	var client *rabbitmq.RabbitMq
	var err error
	client, err = boot.GetRabbitmqClient()

	if err != nil {
		send.Send("", "上传s3错误情况处理失败，rabbitmq 链接失败")
		time.Sleep(10 * time.Second)
		return
	}

	defer client.Close()

	for {

		num := boot.RedisCommonClient.GetClient().SCard(s3key).Val()
		if num == 0 {
			logger.Logger.Info("current time has no s3 upload error")
			time.Sleep(10 * time.Second)
			break
		}
		data, err := boot.RedisCommonClient.GetClient().SPop(s3key).Bytes()
		if err != nil {
			if err == redis.Nil {
				logger.Logger.Info(map[string]interface{}{
					"flag": "s3 error upload task all batch done",
					"err":  err,
				})
			} else {

				logger.Logger.Error(map[string]interface{}{
					"flag": "s3 error upload redis.Client.RPop batch error",
					"err":  err,
				})
			}
			time.Sleep(10 * time.Second)
			return
		}

		var s3data *bean.UploadS3Data
		if err = json.Unmarshal(data, &s3data); err != nil {
			logger.Logger.Error(map[string]interface{}{
				"flag": "s3 error upload json error",
				"err":  err,
			})

			time.Sleep(10 * time.Second)

		} else {

			if err != nil {
				//如果出错 ，将数据还是存入redis里面
				save_file.Saveerrordata(s3data)
				send.Send("", "上传s3错误情况处理失败，rabbitmq 链接失败")
				time.Sleep(10 * time.Second)
			} else {

				if err = client.Send(vars.UploadS3, s3data); err != nil {
					//如果出错 ，将数据还是存入redis里面
					save_file.Saveerrordata(s3data)
					send.Send("", "上传s3错误情况处理失败，发送rabbitmq失败")
					time.Sleep(10 * time.Second)
				}
			}
		}
	}
}
