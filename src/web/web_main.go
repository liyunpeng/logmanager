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

	/*
	app是iri的log级别， 如果不设置log级别， 不会有任何iris的log的输出
	iris log有很多有用的提示， 如http路径对应哪个控制器函数， 如：
	[DBUG] 2019/12/01 16:42 GET: /etcdmanager/add -> controllers.EtcdManangerController.GetAdd()
	[DBUG] 2019/12/01 16:42 GET: /etcdmanager -> controllers.EtcdManangerController.Get()
	[DBUG] 2019/12/01 16:42 Application: 1 registered view engine(s)
	[DBUG] 2019/12/01 16:42 Application: running using 1 host(s)
	[DBUG] 2019/12/01 16:42 Host: addr is localhost:9080
	[DBUG] 2019/12/01 16:42 Host: virtual host is localhost:9080
	[DBUG] 2019/12/01 16:42 Host: register startup notifier
	[DBUG] 2019/12/01 16:42 Host: register server shutdown on interrupt(CTRL+C/CMD+C)
	[DBUG] 2019/12/01 16:42 Host: server will ignore the following errors: [http: Server closed]
	*/
	app.Logger().SetLevel("debug")

	/*
		不管当前代码路径在什么地方， iris.HTML必须基于项目的根目录， 所以是./src/web/views/
	 */
	app.RegisterView(iris.HTML("./src/web/views/", ".html"))

	initConf()

	etcdService := services.NewEtcdService(
		[]string{"127.0.0.1:2379"}, 5 * time.Second)

	etcdKeys := conf.AppConf.GetEtcdKeys()
	for _, key := range etcdKeys {
		resp := etcdService.Get(key)
		for _, ev := range resp.Kvs {
			// return result is not string
			services.ConfChan <- string(ev.Value)
			fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
		}
	}

	go etcdService.EtcdWatch(etcdKeys)

	tailService := services.NewTailService()
	go tailService.RunServer()

	services.NewKafkaService(
		conf.AppConf.KafkaAddr, conf.AppConf.ThreadNum)

	etcdManagerApp := mvc.New(app.Party("/etcdmanager"))
	etcdManagerApp.Register(etcdService)
	etcdManagerApp.Handle(new(controllers.EtcdManangerController))

	fmt.Println("app.run")
	app.Run(
		iris.Addr("localhost:9080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithConfiguration(iris.Configuration{ //默认配置:
			DisableStartupLog:                 false,
			DisableInterruptHandler:           false,
			DisablePathCorrection:             false,
			EnablePathEscape:                  false,
			FireMethodNotAllowed:  		       false,
			DisableBodyConsumptionOnUnmarshal: false,
			DisableAutoFireStatusCode:         false,
			TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
			Charset:                           "UTF-8",
		}),
	)


}
