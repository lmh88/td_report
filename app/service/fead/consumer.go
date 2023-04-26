package fead

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gogf/gf/frame/g"
	"net/http"
	"strings"
	"sync"
	"td_report/app/bean"
	"td_report/common/amazon/sqs_client"
	"td_report/common/kafkaV1"
	"td_report/pkg/logger"
	"td_report/pkg/requests"
	"time"
)

var wg sync.WaitGroup

func Receive(queue string) error {
	logger.Init("fead_" + queue, true)
	logger.Logger.Info(queue, ":start")
	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s", queue))
	client, err := sqs_client.NewSqsClient()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err.Error())
		return err
	}
	//fmt.Println("connect sqs")

	//TODO 测试用gary-sp-traffic
	//kafkaTopic := "amazon_streamapi_sp_ods_traffic"
	kafkaTopic := getKafkaTopic(queue)
	kafkaTool := kafkaV1.NewKafkaTool(g.Cfg().GetString("kafka.url"), kafkaV1.DefaultConfig(), kafkaTopic)
	producer, err := kafkaTool.MyAsyncProducer()
	msg := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Key:   nil,
	}
	//fmt.Println("connect kafka")

	//todo 定时退出
	//ticker := time.NewTicker(time.Second * 30)
	//go func () {
	//	select {
	//	case <- ticker.C:
	//		panic("到时退出")
	//	}
	//} ()

	//todo 开启的协程数
	for i := 0; i < getGoNum(queue); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				err, _ = client.Receivelp(ctx, &queue, 10, func (data *sqs.ReceiveMessageOutput) error {

					if len(data.Messages) > 0 {
						//log.Println(len(data.Messages))
						//records := make([]model.SpTraffic, 0)
						for n, item := range data.Messages {
							body := *item.Body

							if strings.Contains(body, "SubscribeURL") {
								var da bean.FeadConfirm
								if err = json.Unmarshal([]byte(body), &da); err != nil {
									logger.Logger.Error(err.Error(), "confirm sp data error ")
									return err
								}
								logger.Logger.Info("confirm_subscribe", da)
								err = Confirm(ctx, da.SubscribeURL)
								if err != nil {
									return err
								}
								err = client.DelMessage(ctx, &queue, data.Messages[n].ReceiptHandle)
								return err
							} else {
								//msgMap := make(map[string]interface{})
								//err = json.Unmarshal([]byte(body), &msgMap)
								//log.Println(msgMap, err)

								//go func(msg *sarama.ProducerMessage, body string, queue string, n int) {
									msg.Value = sarama.ByteEncoder(body)
									producer.Input() <- msg
									//fmt.Println("push_msg-", queue, n)
									//logger.Logger.Info("push_msg-", queue, n)
									select {
									case suc := <-producer.Successes():
										logger.Logger.Printf("queue:%s, offset: %d", kafkaTopic, suc.Offset)
										//log.Println(kafkaTopic, suc.Offset)

										err := client.DelMessage(ctx, &queue, data.Messages[n].ReceiptHandle)
										if err != nil {
											fmt.Println("del message fail")
											logger.Logger.ErrorWithContext(ctx, err.Error(), "del message fail")
										}
									case fail := <-producer.Errors():
										//fmt.Println("失败")
										logger.Logger.ErrorWithContext(ctx, fail.Err.Error())
										logger.Logger.Println(fail.Err.Error(), "==================fail error============")

									case <-time.After(10 * time.Second):
										// 如果超时，则退出
										logger.Logger.Println("超时退出")

									}
								//}(msg, body, queue, n)

								//tmp := model.SpTraffic{
								//	IdempotencyId: msgMap["idempotency_id"].(string),
								//	CreateDate: gtime.Now(),
								//}
								//records = append(records, tmp)
							}
						}
						//_, err = dao.SpTraffic.Data(records).Insert()
						//if err != nil {
						//	logger.Logger.ErrorWithContext(ctx, err)
						//}
					} else {
						time.Sleep(time.Second * 5)
						fmt.Println(fmt.Sprintf("go:%d:wait 5s"))
					}
					return nil
				})
				if err != nil {
					logger.Logger.ErrorWithContext(ctx, err)
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func Confirm(ctx context.Context, url string) error {
	resp, err := requests.Get(url, requests.WithTimeout(time.Second*60))
	if err != nil {
		fmt.Println(err)
	} else {

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			fmt.Println("success confirm")
			//fmt.Println(string(resp.Body))
			logger.Logger.InfoWithContext(ctx, "success_confirm", string(resp.Body))
		} else {
			fmt.Println(string(resp.Body))
			fmt.Println("error confirm")
			logger.Logger.InfoWithContext(ctx, string(resp.Body), "error_confirm")
		}
	}
	return err
}
