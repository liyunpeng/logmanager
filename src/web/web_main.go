package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"src/services"
	"src/conf"
	"src/web/controllers"
	"fmt"
	"time"
)

func initConf(){
	confFile := "./conf/app.cfg"
	fmt.Println("main")
	err := conf.AppConf.InitConfig(confFile)
	if err != nil {
		fmt.Printf("从本地配置文件加载配置失败:%v", err)
		return
	}

	err = conf.AppConf.InitLogs()

	if err != nil {
		fmt.Printf("init log failed:%v", err)
		return
	}
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

	log的级别从低到高以为	:DEBUG、INFO、WARN、ERROR、FATAL
	只有高于锁设置的级别的以上的log才能打印，
	比如级别设置为info, 那么debug的级别是不打印的
	*/
	app.Logger().SetLevel("debug")

	/*
		不管当前代码路径在什么地方， iris.HTML必须基于项目的根目录， 所以是./src/web/views/
	 */
	app.RegisterView(iris.HTML("./src/web/views/", ".html"))

	initConf()

	etcdService := services.NewEtcdService(
		[]string{"192.168.43.144:2379"}, 5 * time.Second)
		//[]string{"127.0.0.1:2379"}, 5 * time.Second)

	etcdKeys := conf.AppConf.GetEtcdKeys()
	fmt.Println("到etcd服务器，按指定的键遍历键值对")
	for _, key := range etcdKeys {
		resp := etcdService.Get(key)
		for _, ev := range resp.Kvs {
			services.ConfChan <- string(ev.Value)
			fmt.Printf("etcdkey = %s \n etcdvalue = %s \n", ev.Key, ev.Value)
		}
	}

	// 启动对etcd的监听服务，有新的键值对会被监听到
	go etcdService.EtcdWatch(etcdKeys)

	tailService := services.NewTailService()
	go tailService.RunServer()

	services.NewKafkaService(
		conf.AppConf.KafkaAddr, conf.AppConf.ThreadNum)

	/*
		创建iris应用的
		app.Party得到一个路由对象， party的参数就是一个路径，整个应有都是在这个路径下，
		mvc.new 由这个路由对象， 创建一个mvc的app对象。
		这个app就可以做很多事情，
		注册服务啊，注册控制器

	 */
	etcdManagerApp := mvc.New(app.Party("/etcdmanager"))
	etcdManagerApp.Register(etcdService)
	etcdManagerApp.Handle(new(controllers.EtcdManangerController))

	//终端不会打印出此log
	app.Logger().Debug("iris启动服务")

	fmt.Println("iris启动服务")
	/*
		因为只注册了etcdmanager， 访问网址是
		http://localhost:9081/etcdmanager
		如果是http://localhost:9081， 浏览器会显示not found
	 */
	app.Run(
		iris.Addr("localhost:9081"),
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
