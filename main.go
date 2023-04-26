package main

import (
	"github.com/gogf/gf/frame/g"
	_ "td_report/boot"
	_ "td_report/router"
	"td_report/startup"
)

func main() {
	go startup.InitWebMq()
	g.Server().Run()
}
