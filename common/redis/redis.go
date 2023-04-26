package redis

import (
	"crypto/tls"
	"fmt"
	rediS "github.com/go-redis/redis"
	"github.com/gogf/gf/frame/g"
	"td_report/common/file"
	"td_report/pkg/logger"
	"time"
)

var (
	Nil = rediS.Nil
)

type Rds struct {
	client   *rediS.Client
	server   string
	password string
	database int
}

// NewRds 指定初始化连接，针对多个redis，非默认的redis
func NewRds(address string) *Rds {
	return getRds(address)
}

func getRds(address string) *Rds {
	options := &rediS.Options{
		Addr:         g.Cfg().GetString(fmt.Sprintf("%s.host", address)),
		Password:     g.Cfg().GetString(fmt.Sprintf("%s.password", address)),
		DB:           g.Cfg().GetInt(fmt.Sprintf("%s.database", address)),
		PoolSize:     g.Cfg().GetInt(fmt.Sprintf("%s.pool_size", address)),
		MinIdleConns: g.Cfg().GetInt(fmt.Sprintf("%s.min_idle_conns", address)),
		IdleTimeout:  time.Duration(g.Cfg().GetInt(fmt.Sprintf("%s.idle_timeout", address))) * time.Second,
	}

	if g.Cfg().GetInt(fmt.Sprintf("%s.tls", address)) == 1 {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	Client := rediS.NewClient(options)
	_, err := Client.Ping().Result()
	if err != nil {
		logger.Error(map[string]interface{}{
			"flag": "redis init",
			"err":  err.Error(),
		})
	}

	return &Rds{
		client:   Client,
		server:   options.Addr,
		password: options.Password,
		database: options.DB,
	}
}

// NewDefaultRds 默认读取配置文件
func NewDefaultRds() *Rds {
	return getRds("redis")
}

func (c *Rds) initRedis(server string, password string, database int, tls int) {
	c.client = c.poolInitRedis(server, password, database, tls)
}

// redis连接池初始化
func (c *Rds) poolInitRedis(server string, password string, database int, tlsconfig int) *rediS.Client {
	options := &rediS.Options{
		// 连接信息
		Network:  "tcp",    // 网络类型，tcp or unix，默认tcp
		Addr:     server,   // 主机名+冒号+端口，默认localhost:6379
		Password: password, // 密码
		DB:       database, // redis数据库index

		// 连接池容量及闲置连接数量
		PoolSize:     100, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,  // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		// 超时
		DialTimeout:  60 * time.Second, // 连接建立超时时间，默认5秒。
		ReadTimeout:  30 * time.Second, // 读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 30 * time.Second, // 写超时，默认等于读超时
		PoolTimeout:  31 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		// 闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        30 * time.Minute, // 闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		// 命令执行失败时的重试策略
		MaxRetries:      3,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		// 可自定义连接函数
		// Dialer: func() (net.Conn, error) {
		//	netDialer := &net.Dialer{
		//		Timeout:   5 * time.Second,
		//		KeepAlive: 5 * time.Minute,
		//	}
		//
		//	return netDialer.Dial("tcp", server)
		// },

		// 钩子函数
		// OnConnect: func(conn *Rds.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		//	fmt.Printf("conn=%v\n", conn)
		//	return nil
		// },
	}

	if tlsconfig == 1 {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	return rediS.NewClient(options)
}

func (c *Rds) GetClient() *rediS.Client {
	return c.client
}

func (c *Rds) Set(key string, value interface{}, expire time.Duration) error {
	s, _ := file.Serialization(value) // 序列化
	err := c.client.Set(key, s, expire).Err()
	if err != nil {
		g.Log().Println(fmt.Sprintf("Set '%s' err", key), err, value)
	}

	return err
}

func (c *Rds) SetEx(key string, value interface{}, expireTime int) {
	s, _ := file.Serialization(value) // 序列化
	err := c.client.Set(key, s, time.Duration(expireTime)*time.Second).Err()
	if err != nil {
		g.Log().Println(fmt.Sprintf("Set '%s' err", key), err, value)
	}
}

func (c *Rds) SetNx(key string, value interface{}, expireTime time.Duration) bool {
	s, _ := file.Serialization(value) // 序列化
	cmd := c.client.SetNX(key, s, expireTime)
	if cmd.Err() != nil {
		g.Log().Println(fmt.Sprintf("Set '%s' err", key), cmd.Err(), value)
	}

	return cmd.Val()
}

func (c *Rds) Get(key string, v interface{}) {
	result, err := c.client.Get(key).Result()
	if err != nil && err.Error() != "Rds: nil" {
		g.Log().Println(fmt.Sprintf("Get '%s' err", key), err)
	}

	if result != "" {
		_ = file.Deserialization([]byte(result), &v) // 反序列化
	}
}

func (c *Rds) HSet(key string, field string, value interface{}) {
	s, _ := file.Serialization(value) // 序列化
	err := c.client.HSet(key, field, s).Err()
	if err != nil {
		g.Log().Println(fmt.Sprintf("Set '%s' err", key), err, value)
	}
}

func (c *Rds) HGet(key string, field string, v interface{}) {
	result, err := c.client.HGet(key, field).Result()
	if err != nil && err.Error() != "Rds: nil" {
		g.Log().Println(fmt.Sprintf("HGet '%s:%s' err", key, field), err)
	}

	if result != "" {
		_ = file.Deserialization([]byte(result), &v) // 反序列化
	}
}

func (c *Rds) HKeys(key string) []string {
	result, err := c.client.HKeys(key).Result()
	if err != nil && err.Error() != "Rds: nil" {
		g.Log().Println(fmt.Sprintf("HKeys '%s' err", key), err)
	}

	return result
}

func (c *Rds) Expire(key string, t time.Duration) error {
	return c.client.PExpire(key, t).Err()
}

func (c *Rds) Incr(key string) error {
	return c.client.Incr(key).Err()
}

func (c *Rds) Decr(key string) error {
	return c.client.Decr(key).Err()
}

func (c *Rds) Del(key string) error {
	return c.client.Del(key).Err()
}

func (c *Rds) SetLock(lockKey string, expireTime time.Duration) {
	for true {
		if c.SetNx(lockKey, "lock", expireTime) {
			break
		}

		// 等待0.1秒
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Rds) SyncLock(fileName string) bool {
	ok, err := c.client.SetNX(WithSyncLockPrefix(fileName), fileName, time.Second*10).Result()
	if err != nil {
		g.Log().Error(map[string]interface{}{
			"flag": "SyncLock error",
			"err":  err.Error(),
		})

		return false
	}

	return ok
}

func (c *Rds) SyncUnlock(fileName string) {
	err := c.client.Del(WithSyncLockPrefix(fileName)).Err()
	if err != nil {
		g.Log().Error(map[string]interface{}{
			"flag": "SyncUnlock error",
			"err":  err.Error(),
		})
	}
}

type PipeHashData struct {
	Key    string                   `json:"key"`
	Field  string                   `json:"field"`
	Status int                      `json:"status"`
	Data   []map[string]interface{} `json:"data"`
}

// PipeHData PipeHashData redis pipeline data
type PipeHData struct {
	Key           string      `json:"key"`
	Id            int64       `json:"id"`
	Field         string      `json:"field"`
	Status        int         `json:"status"`
	ServingStatus string      `json:"servingStatus"`
	Data          interface{} `json:"data"`
}

// PipeData redis pipeline string data
type PipeData struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Status int    `json:"status"`
}

// Pipeline pipeline
func (c *Rds) Pipeline(key string, field string, data []map[string]interface{}, funcs func(pipe rediS.Pipeliner, key string, field string, data interface{})) error {
	pipe := c.client.Pipeline()
	for _, value := range data {
		funcs(pipe, key, field, value)
	}
	_, err := pipe.Exec()
	return err
}

// HPipeline hash pipeline
func (c *Rds) HPipeline(key string, field string, data []map[string]interface{}) error {
	err := c.Pipeline(key, field, data, func(pipe rediS.Pipeliner, key string, field string, value interface{}) {
		v, _ := file.Serialization(value)
		mp := value.(map[string]interface{})
		fieldValue := mp[field]
		mhfield := fmt.Sprintf("%v", fieldValue)
		pipe.HSet(key, mhfield, v)
	})

	return err
}
