package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/guuid"
	"github.com/google/wire"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"td_report/app/bean"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/common/amazon/sqs_client"
	"td_report/common/kafkaV1"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/report"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

var FeadServiceSet = wire.NewSet(wire.Struct(new(FeadService), "*"))

type SqsCreateSub struct {
	Type             string    `json:"Type"`
	MessageId        string    `json:"MessageId"`
	Token            string    `json:"Token"`
	TopicArn         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	SubscribeURL     string    `json:"SubscribeURL"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
}

type FeadService struct {
	SellerProfileRepository *repo.SellerProfileRepository
	SpConversionRepository  *repo.SpConversionRepository
	SpTrafficRepository     *repo.SpTrafficRepository
	FeadRepository          *repo.FeadRepository
}

const (
	SpConversion = "sp-conversion"
	SpTraffic    = "sp-traffic"
	BudgetUsage  = "budget-usage"
)

var QueueMap = map[string]string{
	SpConversion: "amazon_streamapi_sp_ods_conversion",
	SpTraffic:    "amazon_streamapi_sp_ods_traffic",
	BudgetUsage:  "amazon_streamapi_sp_ods_budget_usage",
}

var FeadMap = map[string]*bean.Postdata{
	SpConversion: &bean.Postdata{
		DataSetId:      "sp-conversion",
		DestinationUri: "arn:aws:sqs:us-east-1:207485092024:sp-conversion",
	},
	SpTraffic: &bean.Postdata{
		DataSetId:      "sp-traffic",
		DestinationUri: "arn:aws:sqs:us-east-1:207485092024:sp-traffic",
	},
	BudgetUsage: &bean.Postdata{
		DataSetId:      "budget-usage",
		DestinationUri: "arn:aws:sqs:us-east-1:207485092024:budget-usage",
	},
}

func NewFeadService(sellerProfileRepository *repo.SellerProfileRepository,
	spConversion *repo.SpConversionRepository,
	spTraffic *repo.SpTrafficRepository,
	feadRepository *repo.FeadRepository,
) *FeadService {
	return &FeadService{
		SellerProfileRepository: sellerProfileRepository,
		SpConversionRepository:  spConversion,
		SpTrafficRepository:     spTraffic,
		FeadRepository:          feadRepository,
	}
}

func (t *FeadService) Fead(clientTag int, profileIdList ...int64) {
	profileTokenClient, err := t.GetProfileTokenList(clientTag, profileIdList...)
	if err != nil {
		logger.Logger.Error(err.Error(), "get data error ")
		return
	}

	if len(profileTokenClient) > 0 {
		for _, item := range profileTokenClient {
			ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%d:%s:%s", item.ProfileId, item.Region, item.ClientId))
			header, err := t.GetconnonHeader(ctx, item)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, err.Error())
				continue
			}

			checkMap := map[string]bool{
				SpTraffic:    false,
				SpConversion: false,
				BudgetUsage:  false,
			}

			err, allList := t.ListAllSubscriptsStream(ctx, item, header)
			if err == nil {
				if len(allList.Subscriptions) > 0 {
					for _, subitem := range allList.Subscriptions {
						checkMap[subitem.DataSetId] = true
					}

					for checkkey, checkitem := range checkMap {
						if checkitem == false {
							sqsCreateSubdata, result := t.CreateStream(ctx, FeadMap[checkkey], item, header)
							logger.Logger.Info("step1:", sqsCreateSubdata)
							if result != nil {
								logger.Logger.ErrorWithContext(ctx, "调创建订阅失败"+checkkey, result.Error())
								continue
							} else {
								logger.Logger.ErrorWithContext(ctx, "调创建订阅成功"+checkkey)
							}

							time.Sleep(2 * time.Second)
						}
					}
				}
			} else {
				logger.Logger.Error(err.Error(), "get data error ")
			}
		}
	} else {
		fmt.Println(" no profile ")
	}
}

func (t *FeadService) GetProfileTokenList(clientTag int, profileIdList ...int64) ([]*bean.ProfileTokenClient, error) {
	return t.SellerProfileRepository.GetProfileAndRefreshToken(clientTag, profileIdList...)
}

func (t *FeadService) GetconnonHeader(ctx context.Context, tokenData *bean.ProfileTokenClient) (map[string]string, error) {
	headers, err := report.GetHeaderWithClient(strconv.FormatInt(tokenData.ProfileId, 10), tokenData.RefreshToken, tokenData.ClientId, tokenData.ClientSecret)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get fead header error")
		return nil, err
	}
	return headers, err
}

// CreateStream 新模式下的stream
func (t *FeadService) CreateStream(ctx context.Context, clientData *bean.Postdata, tokenData *bean.ProfileTokenClient, headers map[string]string) (map[string]interface{}, error) {
	profileIdstr := fmt.Sprintf("%d", tokenData.ProfileId)
	headers["Content-Type"] = "application/vnd.MarketingStreamSubscriptions.StreamSubscriptionResource.v1.0+json"
	region := tokenData.Region
	endpoint := region2.ApiUrl[region]
	path := "/streams/subscriptions"
	url := endpoint + path

	params := map[string]interface{}{
		"dataSetId":          clientData.DataSetId,      //"sp-conversion",
		"destinationArn":     clientData.DestinationUri, //"arn:aws:sqs_client:us-east-1:207485092024:sp-conversion",
		"clientRequestToken": guuid.New(),
		"notes":              profileIdstr,
	}

	data := &model.Fead{
		ProfileId:  tokenData.ProfileId,
		DatasetId:  clientData.DataSetId,
		CreateDate: gtime.Now(),
		Status:     0,
		Step:       1,
	}

	resp, err := requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*60))
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "post  data get error")
		data.StepStatus = 0
		data.ErrReason = err.Error()
		t.FeadRepository.AddOne(data)
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		logger.Logger.ErrorWithContext(ctx, err, "request too many times")
		data.StepStatus = 0
		data.ErrReason = err.Error()
		t.FeadRepository.AddOne(data)
		return nil, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		str := fmt.Sprintf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
		logger.Logger.ErrorWithContext(ctx, str)
		data.StepStatus = 0
		data.ErrReason = str
		t.FeadRepository.AddOne(data)
		return nil, errors.New(str)
	}

	val, err := resp.JsonAndValueIsAny()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get post data json error")
		data.StepStatus = 0
		data.ErrReason = err.Error()
		t.FeadRepository.AddOne(data)
		return nil, errors.New("json get error")
	} else {

		logger.Logger.InfoWithContext(ctx, "success", val)
		fmt.Println(val)
		data.StepStatus = 1
		data.ClientTokenId = val["clientRequestToken"].(string)
		data.MessagesSubscriptionId = val["subscriptionId"].(string)
		t.FeadRepository.AddOne(data)
		return val, nil
	}
}

// Create 调用中的第1
func (t *FeadService) Create(ctx context.Context, clientData *bean.Postdata, tokenData *bean.ProfileTokenClient) (map[string]interface{}, error) {
	headers, err := report.GetHeaderWithClient(strconv.FormatInt(tokenData.ProfileId, 10), tokenData.RefreshToken, tokenData.ClientId, tokenData.ClientSecret)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get fead header error")
		return nil, err
	}

	headers["Content-Type"] = "application/vnd.FEADSubscriptionManagerAPI.MessagesSubscriptionResource.v1.0+json"
	region := tokenData.Region
	endpoint := region2.ApiUrl[region]
	path := "/fead/subscriptions/messages"
	url := endpoint + path
	profileIdstr := fmt.Sprintf("%d", tokenData.ProfileId)
	params := map[string]interface{}{
		"dataSetId":      clientData.DataSetId,      //"sp-conversion",
		"destinationUri": clientData.DestinationUri, //"arn:aws:sqs_client:us-east-1:207485092024:sp-conversion",
		"clientToken":    guuid.New(),
		"notes":          profileIdstr,
	}
	data := &model.Fead{
		ProfileId:  tokenData.ProfileId,
		DatasetId:  clientData.DataSetId,
		CreateDate: gtime.Now(),
		Status:     0,
		Step:       1,
	}

	resp, err := requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*60))
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "post  data get error")
		data.StepStatus = 0
		data.ErrReason = err.Error()
		t.FeadRepository.AddOne(data)
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		logger.Logger.ErrorWithContext(ctx, err, "request too many times")
		data.StepStatus = 0
		data.ErrReason = err.Error()
		t.FeadRepository.AddOne(data)
		return nil, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		str := fmt.Sprintf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
		logger.Logger.ErrorWithContext(ctx, str)
		data.StepStatus = 0
		data.ErrReason = str
		t.FeadRepository.AddOne(data)
		return nil, errors.New(str)
	}

	val, err := resp.JsonAndValueIsAny()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get post data json error")
		data.StepStatus = 0
		data.ErrReason = err.Error()
		t.FeadRepository.AddOne(data)
		return nil, errors.New("json get error")
	} else {

		logger.Logger.InfoWithContext(ctx, "success", val)
		fmt.Println(val)
		data.StepStatus = 1
		data.ClientTokenId = val["clientToken"].(string)
		data.MessagesSubscriptionId = val["messagesSubscriptionId"].(string)
		data.Version = val["version"].(float64)
		t.FeadRepository.AddOne(data)
		return val, nil
	}
}

// Getstatus 根据id查询状态
func (t *FeadService) Getstatus(ctx context.Context, substrId string, tokenData *bean.ProfileTokenClient) {
	path := fmt.Sprintf("/fead/subscriptions/messages/%s", substrId)
	headers, err := report.GetHeaderWithClient(strconv.FormatInt(tokenData.ProfileId, 10), tokenData.RefreshToken, tokenData.ClientId, tokenData.ClientSecret)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get fead header error")
		return
	}

	region := tokenData.Region
	endpoint := region2.ApiUrl[region]

	url := endpoint + path
	fmt.Println(headers)
	fmt.Println(url)
	resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second*60))
	if err != nil {
		fmt.Println(err, "get data error")
	} else {

		fmt.Println(resp.StatusCode, "code")
		fmt.Println(string(resp.Body), "body")
		val, err := resp.JsonAndValueIsAny()
		if err != nil {
			fmt.Println("hahahahhaha2")
		} else {
			fmt.Println("success2==============")
			fmt.Println(val)
		}
	}
}

func (t *FeadService) ListAllSubscriptsStream(ctx context.Context, tokenData *bean.ProfileTokenClient, headers map[string]string) (error, *bean.ListAllStream) {
	region := tokenData.Region
	endpoint := region2.ApiUrl[region]
	//暂时定义为100调数据，不翻页

	url := fmt.Sprintf("%s/streams/subscriptions?maxResults=%d", endpoint, 100)
	resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second*60))
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get data error")
		return err, nil
	} else {

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			var alllist bean.ListAllStream
			err = json.Unmarshal(resp.Body, &alllist)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, err, "get  data json error")
				fmt.Println("查询所有的订阅，json解析失败")
				return errors.New("errors json error"), nil
			} else {
				return nil, &alllist
			}
		} else {

			fmt.Println(resp.StatusCode, "code")
			fmt.Println(string(resp.Body), "body")
			logger.Logger.ErrorWithContext(ctx, string(resp.Body), "get  data error")
			return errors.New("errors"), nil
		}
	}
}

// ListAllSubscripts 列出所有的订阅
func (t *FeadService) ListAllSubscripts(ctx context.Context, tokenData *bean.ProfileTokenClient, headers map[string]string) {
	region := tokenData.Region
	endpoint := region2.ApiUrl[region]
	//暂时定义为100调数据，不翻页
	url := fmt.Sprintf("%s/fead/subscriptions/messages?maxResults=%d", endpoint, 100)
	headers, err := report.GetHeaderWithClient(strconv.FormatInt(tokenData.ProfileId, 10), tokenData.RefreshToken, tokenData.ClientId, tokenData.ClientSecret)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get fead header error")
		return
	}

	resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second*60))
	if err != nil {
		fmt.Println(err, "get data error")
	} else {

		if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusAccepted {
			fmt.Println(resp.StatusCode, "code")
			fmt.Println(string(resp.Body), "body")
			logger.Logger.ErrorWithContext(ctx, err, "get  data error")
		} else {

			fmt.Println(resp.StatusCode, "code")
			fmt.Println(string(resp.Body), "body")
			val, err := resp.JsonAndValueIsAny()
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, err, "get  data json error")
				fmt.Println("查询所有的订阅，json解析失败")
			} else {
				fmt.Println(val)
			}
		}

	}
}

func (t *FeadService) ArchiedStream(ctx context.Context, subId string, tokenData *bean.ProfileTokenClient, headers map[string]string) (bool, error) {
	region := tokenData.Region
	endpoint := region2.ApiUrl[region]
	url := fmt.Sprintf("%s/streams/subscriptions/%s", endpoint, subId)
	params := map[string]interface{}{
		"status": "ARCHIVED",
	}
	headers["Content-Type"] = "application/vnd.MarketingStreamSubscriptions.StreamSubscriptionResource.v1.0+json"
	resp, err := requests.Put(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*60))
	if err != nil {
		fmt.Println(err, "get data error")
		logger.Logger.ErrorWithContext(ctx, err.Error(), "archivd error ")
		return false, err
	} else {
		fmt.Println(string(resp.Body), resp.StatusCode)
		return true, nil
	}
}

// Archied 设置归档，用于将之前的队列归档
func (t *FeadService) Archied(ctx context.Context, subId string, version int, tokenData *bean.ProfileTokenClient) {
	region := tokenData.Region
	endpoint := region2.ApiUrl[region]
	url := fmt.Sprintf("%s/fead/subscriptions/messages/%s", endpoint, subId)
	headers, err := report.GetHeaderWithClient(strconv.FormatInt(tokenData.ProfileId, 10), tokenData.RefreshToken, tokenData.ClientId, tokenData.ClientSecret)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err, "get fead header error")
		return
	}

	params := map[string]interface{}{
		"status":  "ARCHIVED",
		"version": version,
	}

	resp, err := requests.Put(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*60))
	if err != nil {
		fmt.Println(err, "get data error")
	} else {
		fmt.Println(string(resp.Body), resp.StatusCode)
	}
}

func (t *FeadService) Confirme(ctx context.Context, url string) error {
	resp, err := requests.Get(url, requests.WithTimeout(time.Second*60))
	if err != nil {
		fmt.Println(err)
	} else {

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			fmt.Println("success confirme")
			fmt.Println(string(resp.Body))
			logger.Logger.InfoWithContext(ctx, string(resp.Body))
		} else {
			fmt.Println(string(resp.Body))
			fmt.Println("error confirme")
			logger.Logger.InfoWithContext(ctx, string(resp.Body), "error confirme")
		}
	}
	return err
}

// SqsConfirme 调用中第三步
func (t *FeadService) SqsConfirme(ctx context.Context, url string, profileId int64, dataSetId string) error {
	resp, err := requests.Get(url, requests.WithTimeout(time.Second*60))
	data := make(map[string]interface{}, 0)
	data["step"] = 3
	data["update_date"] = gtime.Now()
	if err != nil {
		fmt.Println(err, "第三步请求失败，get data error")
		logger.Logger.ErrorWithContext(ctx, "第三步确认请求失败0", err.Error())
		data["step_status"] = 0
		data["err_reason"] = err.Error()
		t.FeadRepository.Update(profileId, dataSetId, data)
		return err
	} else {

		if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusAccepted {
			str := fmt.Sprintf("code:%d,body:%s", resp.StatusCode, string(resp.Body))
			logger.Logger.ErrorWithContext(ctx, "第三步确认请求失败", str)
			data["step_status"] = 0
			data["err_reason"] = str
			t.FeadRepository.Update(profileId, dataSetId, data)
			return errors.New(str)
		} else {

			val, err := resp.JsonAndValueIsAny()
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, err, "第三步确认解析json失败:get  data json error")
				fmt.Println("查询所有的订阅，json解析失败")
				data["step_status"] = 0
				data["err_reason"] = "第三步确认解析json失败:get  data json error"
				t.FeadRepository.Update(profileId, dataSetId, data)
				return err
			} else {

				str := fmt.Sprintf("code:%d,body:%s", resp.StatusCode, string(resp.Body))
				fmt.Println(val)
				fmt.Println("第三步请求结果成功：", str)
				data["step_status"] = 1
				t.FeadRepository.Update(profileId, dataSetId, data)
				return nil
			}
		}
	}
}

func (t *FeadService) AddSpConversion(body string) error {
	var da model.SpConversion
	if err := json.Unmarshal([]byte(body), &da); err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	da.CreateDate = gtime.Now()
	return t.SpConversionRepository.AddMutils(da)
}

func (t *FeadService) AddSpTraffic(body string) error {
	var da model.SpTraffic
	if err := json.Unmarshal([]byte(body), &da); err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	da.CreateDate = gtime.Now()
	return t.SpTrafficRepository.AddMutils(da)
}

func (t *FeadService) BudgetUsage(body string) error {
	return nil
}

var queueList = []string{SpConversion, SpTraffic, BudgetUsage}

var stopLabel bool = false

// Receivelp 长链接处理消息，提高消息消费能力
func (t *FeadService) Receivelp(topic string) error {
	var flag = false
	for _, item := range queueList {
		if item == topic {
			flag = true
			break
		}
	}

	if flag == false {
		return errors.New("paramas error the queue name is not exists ")
	}

	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s", topic))
	client, err := sqs_client.NewSqsClient()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err.Error())
		return err
	}

	kafkatopic := QueueMap[topic]
	kafkaTool := kafkaV1.NewKafkaTool(g.Cfg().GetString("kafka.url"), kafkaV1.DefaultConfig(), kafkatopic)
	producer, err := kafkaTool.MyAsyncProducer()
	msg := &sarama.ProducerMessage{
		Topic: kafkatopic,
		Key:   nil,
	}

	lock := sync.Mutex{}
	num := g.Cfg().GetInt32("sqs.consumer_num")
	for {

		if stopLabel {
			fmt.Println("===========休眠1min退出循环===============")
			time.Sleep(1 * time.Minute)
			break
		}

		if err, result := client.Receivelp(ctx, &topic, num, func(data *sqs.ReceiveMessageOutput) error {
			for _, item := range (*data).Messages {
				body := *item.Body
				if strings.Contains(body, "SubscribeURL") {
					var da bean.FeadConfirm
					if err := json.Unmarshal([]byte(body), &da); err != nil {
						logger.Logger.Error(err.Error(), "confirm data error ")
						return err
					}

					//中间关联关系断开了,如果是多个profileid 目前无法判断是哪一个订阅了哪一个没有
					return t.Confirme(ctx, da.SubscribeURL)

				} else {

					go func(mymsg *sarama.ProducerMessage, mybody string, mytopic string) {
						lock.Lock()
						defer lock.Unlock()
						mymsg.Value = sarama.ByteEncoder(mybody)
						producer.Input() <- mymsg
						select {
						case suc := <-producer.Successes():
							logger.Logger.Printf("topic:%s, offset: %d,  timestamp: %s\n", mytopic, suc.Offset, time.Now().Format(vars.TIMEFORMAT))
							err = client.DelMessage(ctx, &mytopic, data.Messages[0].ReceiptHandle)
							if err != nil {
								logger.Logger.ErrorWithContext(ctx, err.Error(), "del message fail")
							}
						case fail := <-producer.Errors():
							logger.Logger.ErrorWithContext(ctx, fail.Err.Error())
							logger.Logger.Info(fail.Err.Error(), "==================fail error============")
							stopLabel = true
						case <-time.After(15 * time.Second):
							// 如果超时，则退出
							logger.Logger.Info("超时退出")
							stopLabel = true
						}
					}(msg, body, topic)
				}
			}

			return nil
		}); err != nil || result == false {
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, err.Error())
				fmt.Println("=====error", err.Error())
			}

			time.Sleep(50 * time.Millisecond)
			stopLabel = true
		}
	}

	return nil
}

func (t *FeadService) Receive(topic string) error {
	var flag = false
	for _, item := range queueList {
		if item == topic {
			flag = true
			break
		}
	}

	if flag == false {
		return errors.New("paramas error the queue name is not exists ")
	}

	ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s", topic))
	client, err := sqs_client.NewSqsClient()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, err.Error())
		return err
	}

	//尝试了全部放开协成的处理方式，内存飙升消息sqs报错，目前只能退而求其次处理的时候单个协成single，并发数达到sqs的警戒线，可以尝试限制线程数来跑10个以下
	single := make(chan struct{}, 10)
	kafkatopic := QueueMap[topic]
	kafkaTool := kafkaV1.NewKafkaTool(g.Cfg().GetString("kafka.url"), kafkaV1.DefaultConfig(), kafkatopic)
	producer, err := kafkaTool.MyAsyncProducer()
	msg := &sarama.ProducerMessage{
		Topic: kafkatopic,
		Key:   nil,
	}

	for {

		if stopLabel {
			fmt.Println("===========休眠1min退出循环===============")
			time.Sleep(1 * time.Minute)
			break
		}

		go func() {
			if err, result := client.Receive(ctx, &topic, 8, func(data *sqs.ReceiveMessageOutput) error {
				if len((*data).Messages) > 0 {
					body := *(*data).Messages[0].Body
					if strings.Contains(body, "SubscribeURL") {
						var da bean.FeadConfirm
						if err := json.Unmarshal([]byte(body), &da); err != nil {
							logger.Logger.Error(err.Error(), "confirm data error ")
							return err
						}

						//中间关联关系断开了,如果是多个profileid 目前无法判断是哪一个订阅了哪一个没有
						return t.Confirme(ctx, da.SubscribeURL)

					} else {

						go func(mymsg *sarama.ProducerMessage, mybody string, mytopic string) {
							mymsg.Value = sarama.ByteEncoder(mybody)
							producer.Input() <- mymsg
							select {
							case suc := <-producer.Successes():
								logger.Logger.Printf("topic:%s, offset: %d,  timestamp: %s\n", mytopic, suc.Offset, time.Now().Format(vars.TIMEFORMAT))
								err = client.DelMessage(ctx, &mytopic, data.Messages[0].ReceiptHandle)
								if err != nil {
									logger.Logger.ErrorWithContext(ctx, err.Error(), "del message fail")
								}
							case fail := <-producer.Errors():
								logger.Logger.ErrorWithContext(ctx, fail.Err.Error())
								logger.Logger.Println(fail.Err.Error(), "==================fail error============")
								stopLabel = true
							case <-time.After(35 * time.Second):
								// 如果超时，则退出
								logger.Logger.Println("超时退出")
								stopLabel = true
							}
						}(msg, body, topic)
					}
				} else {
					fmt.Println("empty queue")
					stopLabel = true
				}

				return nil
			}); err != nil || result == false {
				if err != nil {
					logger.Logger.ErrorWithContext(ctx, err.Error())
					fmt.Println("=====error", err.Error())
				}

				time.Sleep(50 * time.Millisecond)
				stopLabel = true
			}
			single <- struct{}{}
		}()

		<-single
	}

	return nil
}
