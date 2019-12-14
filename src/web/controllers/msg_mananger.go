package controllers

//import (
//	//"src/conf"
//	//"src/services"
//	//"fmt"
//	"github.com/kataras/iris"
//	"services"
//)
//
//type MsgManangerController struct {
//	ctx iris.Context
//
//	etcdService services.EtcdService
//}
//
//func (c *MsgManangerController) Get() {
//
//	//etcdKeys := conf.AppConf.GetEtcdKeys()
//	//for _, key := range etcdKeys {
//	//	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	//
//	//	resp := c.etcdService.Get(key)
//	//
//	//	for _, ev := range resp.Kvs {
//	//		// return result is not string
//	//		services.ConfChan <- string(ev.Value)
//	//		fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
//	//	}
//	//}
//}
