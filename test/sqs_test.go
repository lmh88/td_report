package test

import (
	"fmt"
	"td_report/common/tool"
	"testing"
)

func TestSqlReceive(t *testing.T) {
	//topic := "sp-conversion"
	//ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s", topic))
	//sqs_client.Receive(ctx, &topic, 1, func(data *sqs.ReceiveMessageOutput) error {
	//	if len((*data).Messages) > 0 {
	//		for _, item := range (*data).Messages {
	//			fmt.Println("Message ID:     " + *item.MessageId)
	//			fmt.Println("Message Handle: " + *item.ReceiptHandle)
	//			fmt.Println("body:", *item.Body)
	//		}
	//	}
	//
	//	return nil
	//})
}

func TestListQueue(t *testing.T) {
	//topic := "sp-conversion"
	//sqs_client.ListQueue(topic)
	token:="Atzr|IwEBIK0MRzpvGx9-oR-cR-JCqiQsuKMEy7eYTFBtOtDOc0BxC97POmBNyNjy57Shvs_waUhm-XM_sCzJ6-xP8ZhI7bDQ7a4XKzy7bOyjVO3RSEEhKW2GPIEaJIUFk1vgitDROCHQ6oonk6_UzqoaTXZdyfSABykhBfBdCwO_kPZy_px2p7tiV1vl0X0rROnzgBbgiGnEik64XRWq-DTr5qI9ay0vRjp99w9CH_IoRT5hI-FSct4x6Z5VfxEutpeRzdHt1RI2ja8829dQe-IWG-jIYnpzNnQa2GauMwaIMCSq8ymkYMdKRrC2da-lwmFZHGO8AM8FVhCNQvL9extJChwFsdw7xHltywK8AUvHQ7WXMo5eqYnhC6k7zQ4kOQgFEzpW0AjrZ7DFHEt-PSXNYoknSI7Kabpv7KNe6gC4LD_uySKjF81BO9W8-antfCBHx_Px_XvMfBWOj1bwIQ-ICubgUrCi"
	fmt.Println(tool.GetMd5(token))
}
