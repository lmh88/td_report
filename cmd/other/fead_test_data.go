package other

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/service/report"
	"td_report/pkg/logger"
)

var step string
var FeadtestCmd = &cobra.Command{
	Use:   "feadtest",
	Short: "feadtest",
	Long:  `feadtest`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("fead_test", false)
		logger.Logger.Info("fead_test called")
		feadService := app.InitializeFeadService()
		profileId := int64(3597281164691753)
		setp2data := map[string]interface{}{
			"clientToken":            "81938b3c-ed2a-42da-8323-8f5f2d0baadc",
			"messagesSubscriptionId": "amzn1.fead.cs1.nGpPMLLjtaSwfUerNSBGDQ",
			"version":                1,
		}

		profileTokenClient, err := feadService.GetProfileTokenList(1, profileId)
		if err != nil || len(profileTokenClient) == 0 {
			fmt.Println("get data error or profile client is empty ")
			return
		}
		step = "5"
		item := profileTokenClient[0]
		ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%d:%s:%s", profileId, item.Region, item.ClientId))
		data := &bean.Postdata{
			DataSetId:      "sp-conversion",
			DestinationUri: "arn:aws:sqs:us-east-1:207485092024:sp-conversion",
		}
		//data := &bean.Postdata{
		//	DataSetId:      "budget-usage",
		//	DestinationUri: "arn:aws:sqs:us-east-1:207485092024:budget-usage",
		//}
		if step == "1" {
			fmt.Println("fead called-1")
			logger.Logger.Info("fead called-1")
			ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%d:%s:%s", item.ProfileId, item.Region, item.ClientId))
			header, _ := feadService.GetconnonHeader(ctx, item)
			retdata, err := feadService.CreateStream(ctx, data, item, header)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				fmt.Println(retdata)
			}
		}

		if step == "2" {
			fmt.Println("fead called-3")
			logger.Logger.Info("fead called-3")
			url := "https://sns.us-east-1.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-east-1:802324068763:FEADMessageDeliveryTopic-sp-conversion-A2EUQ1WTGCTBG2-ENTITY2NARS0YSTL83U&Token=2336412f37fb687f5d51e6e2425dacbba931e5f28828ef60ce42fdaaec969a60ee799a997ed1251fea9b74d410054a7c6a981ef6428ad9fa9bfd931162f991a36fc493ba4e6e5664aa57420d57e6fb8a92bc2f8859b2d4cc68bedc126313259955d3e5f4646399b77a0aeb340db41c083f53cb9e5ff7473e1b366475ce65e17d8aa4fb01bf3f9839d795dfa5ac23cdad96e637a2540a103c5f704e6a8cf26a8933c5127de5fbaa7f3914ebb4f03f7eb5"
			Confirm(profileId, feadService, url, data.DataSetId)
		}

		// 查看状态
		if step == "3" {
			feadService.Getstatus(ctx, setp2data["messagesSubscriptionId"].(string), item)
		}

		// 归档
		if step == "4" {
			header, _ := feadService.GetconnonHeader(ctx, item)
			feadService.ArchiedStream(ctx, setp2data["messagesSubscriptionId"].(string), item, header)
		}
		// 查看多个订阅状态
		if step == "5" {
			header, _ := feadService.GetconnonHeader(ctx, item)
			feadService.ListAllSubscriptsStream(ctx, item, header)
		}
	},
}

func init() {
	RootCmd.AddCommand(FeadtestCmd)
}

func Confirm(profileId int64, feadService *report.FeadService, url string, datasetId string) {
	profileTokenClient, err := feadService.GetProfileTokenList(1, profileId)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(profileTokenClient) == 0 {
		fmt.Println("数据错误")
		return
	}

	item := profileTokenClient[0]
	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%d:%s:%s", item.ProfileId, item.Region, item.ClientId))
	feadService.SqsConfirme(ctx, url, profileId, datasetId)
}
