package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/bean"
	"td_report/common/amazon/s3"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"
)

var UploadS3taskCmd = &cobra.Command{
	Use:   "uploads3_task",
	Short: "uploads3_task",
	Long:  `uploads3_task`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("uploads3_task", false)
		logger.Logger.Info("uploads3_task called", time.Now().Format(vars.TIMEFORMAT))
		uploads3taskCmdfunc()
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func uploads3taskCmdfunc() {
	rabbitMq, err := rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	if err != nil {
		logger.Logger.Error(err)
		return
	}

	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {
		logger.Logger.Error(err)
		return
	}

	uploadIsOpen := g.Cfg().GetBool("server.S3Upload")
	if uploadIsOpen == true {
		rabbitMq.ReceivWithAck(vars.UploadS3, func(bytes []byte) error {
			var uploaddata *bean.UploadS3Data
			if err = json.Unmarshal(bytes, &uploaddata); err != nil {
				logger.Logger.Error(err)
				return err
			}

			var ctx = logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s", uploaddata.Key))
			if err = client.UploadSingle(ctx, uploaddata.Path, uploaddata.Key); err != nil {
				save_file.Saveerrordata(uploaddata)
				return err
			} else {
				return nil
			}
		})

	} else {

		rabbitMq.ReceivWithAck(vars.UploadS3, func(bytes []byte) error {
			var uploaddata *bean.UploadS3Data
			if err = json.Unmarshal(bytes, &uploaddata); err != nil {
				logger.Logger.Error(err)
				return err
			}

			fmt.Println(uploaddata.Key, uploaddata.Path)
			return nil
		})
	}

}
