package boot

import (
	"github.com/gogf/gf/frame/g"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/common/redis"
	"td_report/common/sendmsg/wechart"
	_ "td_report/packed"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"
	"time"
)

var Redisclient *redis.Rds
var RedisCommonClient *redis.Rds
var Rlimiter *limiter.Limiter
var RabitmqClient *rabbitmq.RabbitMq

func init() {
	Redisclient = redis.NewDefaultRds()
	RedisCommonClient = redis.NewRds("redis_common")
	Rlimiter = limiter.NewLimiter(Redisclient.GetClient())
	var err error
	RabitmqClient, err = rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	if err != nil {
		logger.Logger.Error(err, "get rabbitmq client is error ")
		panic(err)
	}
}

func GetRabbitmqClient() (*rabbitmq.RabbitMq, error) {
	var (
		err    error
		n      = 5
		client *rabbitmq.RabbitMq
	)

	//如果rabbitmq 连接失败重连5次
	client, err = rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	if err != nil {
		for i := 0; i < n; i++ {
			client, err = rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
			if err == nil {
				break
			}

			time.Sleep(2 * time.Second)
		}

		if err != nil || client == nil {
			key := g.Cfg().GetString("wechat.key")
			send := wechart.NewSendMsg(key, g.Cfg().GetBool("wechat.open"))
			send.Send("", "rabbit mq connect error !!!")
		}
	}

	return client, err
}
