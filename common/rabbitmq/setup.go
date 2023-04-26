package rabbitmq

import "github.com/gogf/gf/frame/g"

var MqServer EventServerIface

func NewMqServer() error {
	address := g.Cfg().GetString("rabbitmq.address")
	config := Config{
		addres:       address,
		BusinessName: "",
	}
	msServer, err := NewEventServer(config)
	if err != nil {
		return err
	}

	err = msServer.Start()
	if err != nil {
		return err
	}
	MqServer = msServer
	return nil
}
