package controllers

import (
	"../../services"
	"github.com/kataras/iris"
	"time"
)

type MsgManangerController struct {
	ctx iris.Context

	etcdService services.EtcdService
}

func (c *MsgManangerController) f() {


	for _, key := range etcdKeys {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		resp, err := cli.Get(ctx, key)

		fmt.Println("resp:", resp)
		cancel()
		if err != nil {
			fmt.Println("get etcd key failed, error:", err)
			continue
		}

		for _, ev := range resp.Kvs {
			// return result is not string
			confChan <- string(ev.Value)
			fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
		}
	}
}
