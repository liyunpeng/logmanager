package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func WebMain() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	users := mvc.New(app.Party("/users"))

	// Bind the "userService" to the UserController's Service (interface) field.
	//users.Register(userService)

	users.Handle(new(controllers.ConfManangerController))


}
