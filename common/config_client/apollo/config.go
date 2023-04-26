package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/gogf/gf/frame/g"
	"github.com/sirupsen/logrus"
)

var appolloConfig *config.AppConfig

func init() {
	//实体机上暂时读取配置文件，docker或者k8s后读取环境变量
	appolloConfig = &config.AppConfig{
		AppID:          g.Cfg().GetString("apollo.AppID"),
		Cluster:        g.Cfg().GetString("apollo.Cluster"),
		IP:             g.Cfg().GetString("apollo.IP"),
		NamespaceName:  g.Cfg().GetString("apollo.NamespaceName"),
		IsBackupConfig: g.Cfg().GetBool("apollo.IsBackupConfig"),
		Secret:         g.Cfg().GetString("apollo.Secret"),
	}
}

func GetClient() (agollo.Client, error) {
	agollo.SetLogger(&logrus.Logger{})
	return agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return appolloConfig, nil
	})
}

func Getconfig(key string) (map[string]interface{}, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	cache := client.GetConfigCache(appolloConfig.NamespaceName)
	if value, err := cache.Get(key); err != nil {
		return nil, err
	} else {
		data := map[string]interface{}{
			key: value,
		}
		return data, nil
	}
}

func GetconfigBykeys(keys []string) (map[string]interface{}, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	cache := client.GetConfigCache(appolloConfig.NamespaceName)
	data := make(map[string]interface{})
	for _, key := range keys {
		if value, err := cache.Get(key); err != nil {
			return nil, err
		} else {
			data[key] = value
		}
	}

	return data, nil
}
