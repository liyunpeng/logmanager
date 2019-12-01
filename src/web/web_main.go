package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"../services"
)

func WebMain() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	users := mvc.New(app.Party("/users"))

	// Bind the "userService" to the UserController's Service (interface) field.

	etcdService := services.NewEtcdService([]string{}, 10)

	users.Register(etcdService)

	users.Handle(new(controllers.ConfManangerController))
}
