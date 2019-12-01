package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

func main() {


	timeout := time.Duration(appConf.EtcdTimeOut) * time.Second
	var etcdAddrSlice []string
	etcdAddrSlice = append(etcdAddrSlice, appConf.EtcdAddr)
	err = initEtcd(etcdAddrSlice, appConf.EtcdWatchKey, timeout)
	if err != nil {
		logs.Error("init etcd Failed:%v", err)
		return
	}
	fmt.Println("init etcd success")

	err = initKafka(appConf.KafkaAddr, appConf.ThreadNum)
	if err != nil {
		logs.Error("init kafka Failed:%v", err)
		return
	}
	fmt.Println("init kafka success")

	runServer()
}
