package controllers

import (
	"github.com/kataras/iris"
	"../../services"
)

type ConfManangerController struct {
	Ctx iris.Context

	Service services.EtcdService
}


func  (c *ConfManangerController)f(){
	etcdKey := "/logagent/192.168.0.142/logconfig"

	etcdValue := `
	[
		{
			"topic":"nginx_log",
			"log_path":"/d/log1",
			"service":"test_service",
			"send_rate":1000
		},
			
		{
			"topic":"nginx_log",
			"log_path":"/d/log1",
			"service":"test_service",
			"send_rate":1000
		}
	]`

	c.Service.PutKV(etcdKey, etcdValue)
}