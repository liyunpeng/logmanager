package controllers

import (
	"github.com/kataras/iris"
	"../../services"
)

type ConfManangerController struct {
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.EtcdService


}