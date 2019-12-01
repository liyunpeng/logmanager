package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"../services"
	"../conf"
	"./controllers"
	"fmt"
	"time"
)

func initConf(){
	confFile := "./conf/app.cfg"
	fmt.Println("main")
	err := conf.AppConf.InitConfig(confFile)
	if err != nil {
		fmt.Printf("init conf failed:%v", err)
		return
	}
	fmt.Println("init conf success")


	err = conf.AppConf.InitLogs()

	if err != nil {
		fmt.Printf("init log failed:%v", err)
		return
	}
	fmt.Println("init logs success")
}

func WebMain() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	app.RegisterView(iris.HTML("./src/web/views/", ".html"))


	initConf()

	etcdService := services.NewEtcdService(
		[]string{"127.0.0.1"}, 5 * time.Second)
	etcdKeys := conf.AppConf.GetEtcdKeys()
	etcdService.EtcdWatch(etcdKeys)

	tailService := services.NewTailService()
	tailService.RunServer()

	services.NewKafkaService(
		conf.AppConf.KafkaAddr, conf.AppConf.ThreadNum)


	etcdManagerApp := mvc.New(app.Party("/etcdmanager"))
	etcdManagerApp.Register(etcdService)
	etcdManagerApp.Handle(new(controllers.EtcdManangerController))
}
