package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"../../services"
	"github.com/kataras/iris/mvc"
)

type EtcdManangerController struct {
	Ctx iris.Context

	Service services.EtcdService
}

var v = mvc.View{
	Name:"conf_manager.html",
}

func (e *EtcdManangerController) Get() mvc.Result{
	return v
}

func  (e *EtcdManangerController)GetAdd(){
	etcdKey := "/logagent/192.168.0.142/logconfig"

	etcdValue := `
	[
		{
			"topic":"nginx_log",
			"log_path":"D:\\log1",
			"service":"test_service",
			"send_rate":1000
		},
			
		{
			"topic":"nginx_log",
			"log_path":"D:\\log1",
			"service":"test_service",
			"send_rate":1000
		}
	]`

	e.Service.PutKV(etcdKey, etcdValue)

	fmt.Println("etcd putkv")
}

//func  (e *EtcdManangerController)PostAdd() mvc.Result{
//	f := e.Ctx.FormValue("data")
//	e.Ctx
//
//	return mvc.Response{
//		//如果不是nil，则会显示此错误
//		Err: err,
//	}
//}