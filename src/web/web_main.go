package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"../services"
	"../datasource"
	"./controllers"
	"fmt"
)

func initConf(){
	confFile := "./conf/app.cfg"
	fmt.Println("main")
	err := datasource.AppConf.InitConfig(confFile)
	if err != nil {
		fmt.Printf("init conf failed:%v", err)
		return
	}
	fmt.Println("init conf success")


	err = datasource.AppConf.InitLogs()

	if err != nil {
		fmt.Printf("init log failed:%v", err)
		return
	}
	fmt.Println("init logs success")
}

func WebMain() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	initConf()

	users := mvc.New(app.Party("/users"))

	// Bind the "userService" to the UserController's Service (interface) field.

	etcdService := services.NewEtcdService([]string{}, 10)

	users.Register(etcdService)

	users.Handle(new(controllers.ConfManangerController))
}
