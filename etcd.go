package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	//client "github.com/coreos/etcd/clientv3"
	client "go.etcd.io/etcd/clientv3"
)

var (
	confChan  = make(chan string, 10)
	cli       *client.Client
	waitGroup sync.WaitGroup
)

func initEtcd(addr []string, keyFormat string, timeout time.Duration) (err error) {
	// init a global var cli and can not close
	fmt.Println("endpoints:", addr)

	cli, err = client.New(client.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})

	kv := client.NewKV(cli)


	if err != nil {
		fmt.Println("connect etcd error:", err)
		return
	}
	logs.Debug("init etcd success")
	// defer cli.Close()   //can not close

	var etcdKeys []string
	ips, err := getLocalIP()
	if err != nil {
		fmt.Println("get local ip error:", err)
		return
	}
	for _, ip := range ips {
		key := fmt.Sprintf(keyFormat, ip)
		etcdKeys = append(etcdKeys, key)
	}
	fmt.Println("etcdkeys:", etcdKeys)

	//kapi := cli.NewKeysAPI(c)

	ctx, _ := context.WithTimeout(context.Background(), 5 *time.Second)

	//putResp,err := kv.Put(ctx,"/job/v3","push the box")  //withPrevKV()是为了获取操作前已经有的key-value
	//if err != nil{
	//	panic(err)
	//}
	//fmt.Printf("kvs1: %v",putResp.PrevKv)

	getResp,err := kv.Get(ctx,"/job/v3") //withPrefix()是未了获取该key为前缀的所有key-value
	if err != nil{
		panic(err)
	}
	fmt.Printf("kvs2:  %v",getResp.Kvs)

	/*
		{
		"service":"test_service",
		"log_path":"/search/nginx/logs/ping-android.shouji.sogoucom_access_log",
	    "topic": "nginx_log",
		"send_rate": 1000
		},

	type LogConfig struct {
		Topic    string `json:"topic"`
		LogPath  string `json:"log_path"`
		Service  string `json:"service"`
		SendRate int    `json:"send_rate"`
	}

	 */

	s := `
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

	s1 := "/logagent/192.168.0.142/logconfig"
	putResp,err := kv.Put(ctx, s1, s)  //withPrevKV()是为了获取操作前已经有的key-value
	if err != nil{
		panic(err)
	}
	fmt.Printf("kvs1: %v",putResp.PrevKv)

	//getResp1 ,err := kv.Get(ctx,s1) //withPrefix()是未了获取该key为前缀的所有key-value
	//if err != nil{
	//	panic(err)
	//}
	//fmt.Printf("kvs2:  %v",getResp1.Kvs)

	// first, pull conf from etcd
	for _, key := range etcdKeys {
		ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)

		//cli.Put(ctx, "sample_key_1", "sample_value")
		//
		//
		//fmt.Println("ectd key :", key)
		//
		//time.Sleep(5*time.Second)
		//
		//resp1, err := cli.Get(ctx, "sample_key_1")
		//
		//fmt.Println("resp1:", resp1)


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

	waitGroup.Add(1)
	// second, start a goroutine to watch etcd
	go etcdWatch(etcdKeys)
	return
}

// watch etcd
func etcdWatch(keys []string) {
	defer waitGroup.Done()

	var watchChans []client.WatchChan
	for _, key := range keys {
		rch := cli.Watch(context.Background(), key)
		fmt.Println("rch:", rch)

		watchChans = append(watchChans, rch)
	}

	for {
		for _, watchC := range watchChans {
			select {
			case wresp := <-watchC:
				for _, ev := range wresp.Events {
					confChan <- string(ev.Kv.Value)
					logs.Debug("etcd key = %s , etcd value = %s", ev.Kv.Key, ev.Kv.Value)
				}
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

//GetEtcdConfChan is func get etcd conf add to chan
func GetEtcdConfChan() chan string {
	return confChan
}
