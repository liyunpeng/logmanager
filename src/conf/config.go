package conf

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var AppConf = &AppConfig{}

type AppConfig struct {
	EtcdAddr     string
	EtcdTimeOut  int
	EtcdWatchKey string

	KafkaAddr string

	ThreadNum int
	LogFile   string
	LogLevel  string
}


func (a *AppConfig )InitConfig(file string) (err error) {
	conf, err := config.NewConfig("ini", file)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	AppConf.EtcdAddr = conf.String("etcd_addr")
	AppConf.EtcdTimeOut = conf.DefaultInt("etcd_timeout", 5)
	AppConf.EtcdWatchKey = conf.String("etcd_watch_key")

	AppConf.KafkaAddr = conf.String("kafka_addr")

	AppConf.ThreadNum = conf.DefaultInt("thread_num", 4)
	AppConf.LogFile = conf.String("log")
	AppConf.LogLevel = conf.String("level")
	fmt.Println("从本地配置文件", file, "加载到的配置信息", AppConf)
	return
}

func (a *AppConfig )InitLogs() (err error) {

	config := make(map[string]interface{})
	config["filename"] = AppConf.LogFile
	config["level"] = getLevel(AppConf.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	// logs.SetLogFuncCall(true) // print file name and row number
	return
}


func getLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "trace":
		return logs.LevelTrace
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "error":
		return logs.LevelError
	default:
		return logs.LevelDebug
	}
}

func (e *AppConfig) GetEtcdKeys() ([]string) {
	var etcdKeys []string
	ips, err := getLocalIP()
	if err != nil {
		fmt.Println("get local ip error:", err)
		//return err
	}
	for _, ip := range ips {
		key := fmt.Sprintf(AppConf.EtcdWatchKey, ip)
		etcdKeys = append(etcdKeys, key)
	}
	fmt.Println("从etcd服务器获取到的以IP名为键的键值对: \n", etcdKeys)
	return etcdKeys
}