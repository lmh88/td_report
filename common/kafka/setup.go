package kafka

var MqServer EventServerIface

func NewMqServer() error {
	var serverPath = "url"
	var address []string = []string {
		serverPath,
	}
	config:= Config{
		addres: address,
		BusinessName: "",
	}
	msServer,err:= NewEventServer(config)
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

