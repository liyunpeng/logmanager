package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"../conf"
)

type TailService interface {
	RunServer()
}

type tailService struct {

}


func (t *tailService ) RunServer() {
	tailMgr = NewTailMgr()
	tailMgr.Process()
	waitGroup.Wait()
}


func NewTailService() *tailService{

	return &tailService{}
}
// TailObj is TailMgr's instance
type TailObj struct {
	tail     *tail.Tail
	offset   int64
	logConf  conf.LogConfig
	secLimit *conf.SecondLimit
	exitChan chan bool
}

var tailMgr *TailMgr

//TailMgr to manage tailObj
type TailMgr struct {
	tailObjMap map[string]*TailObj
	lock       sync.Mutex
}

// NewTailMgr init TailMgr obj
func NewTailMgr() *TailMgr {
	return &TailMgr{
		tailObjMap: make(map[string]*TailObj, 16),
	}
}

//AddLogFile to Add tail obj
func (t *TailMgr) AddLogFile(logConfig conf.LogConfig) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	fmt.Println("add log file:", logConfig)
	_, ok := t.tailObjMap[logConfig.LogPath]
	if ok {
		err = fmt.Errorf("duplicate filename:%s", logConfig.LogPath)
		return
	}

	tail, err := tail.TailFile(logConfig.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // read to tail
		MustExist: false,                                //file does not exist, it does not return an error
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}

	tailObj := &TailObj{
		tail:     tail,
		offset:   0,
		logConf:  logConfig,
		secLimit: conf.NewSecondLimit(int32(logConfig.SendRate)),
		exitChan: make(chan bool, 1),
	}
	t.tailObjMap[logConfig.LogPath] = tailObj

	waitGroup.Add(1)
	go tailObj.readLog()
	return
}

func (t *TailMgr) reloadConfig(logConfArr []conf.LogConfig) (err error) {
	for _, logConfArrValue := range logConfArr {
		tailObj, ok := t.tailObjMap[logConfArrValue.LogPath]

		if !ok {
			err = t.AddLogFile(logConfArrValue)
			if err != nil {
				logs.Error("add log file failed:%v", err)
				continue
			}
			continue
		}
		tailObj.logConf = logConfArrValue
		tailObj.secLimit.Limit = int32(logConfArrValue.SendRate)
		t.tailObjMap[logConfArrValue.LogPath] = tailObj
		fmt.Println("tailObj:", tailObj)
	}

	for key, tailObj := range t.tailObjMap {
		var found = false
		for _, newValue := range logConfArr {
			if key == newValue.LogPath {
				found = true
				break
			}
		}
		if found == false {
			logs.Warn("log path :%s is remove", key)
			tailObj.exitChan <- true
			delete(t.tailObjMap, key)
		}
	}
	return
}

// Process hava two func get new log conf and reload conf
func (t *TailMgr) Process() {
	for etcdConfValue := range ConfChan {
		logs.Debug("log etcdConfValue: %v", etcdConfValue)
		//fmt.Printf("log etcdConfValue: %v", etcdConfValue)

		var logConfArr []conf.LogConfig

		err := json.Unmarshal([]byte(etcdConfValue), &logConfArr)
		fmt.Println("logConfArr: ", logConfArr)


		if err != nil {
			logs.Error("unmarshal failed, err: %v etcdConfValue :%s", err, etcdConfValue)
			fmt.Println("unmarshal failed, err: %v etcdConfValue :%s", err, etcdConfValue)
			continue
		}

		err = t.reloadConfig(logConfArr)
		if err != nil {
			logs.Error("reload config from etcd failed: %v", err)
			continue
		}
		//logs.Debug("reload config from etcd success")
		fmt.Printf("reload config from etcd success")
	}
}

func (t *TailObj) readLog() {

	for line := range t.tail.Lines {
		if line.Err != nil {
			logs.Error("read line error:%v ", line.Err)
			continue
		}

		lineStr := strings.TrimSpace(line.Text)
		fmt.Println("readLog :", lineStr)

		if len(lineStr) == 0 || lineStr[0] == '\n' {
			continue
		}

		kafkaSender.addMessage(line.Text, t.logConf.Topic)
		t.secLimit.Add(1)
		t.secLimit.Wait()

		select {
		case <-t.exitChan:
			logs.Warn("tail obj is exited: config:", t.logConf)
			return
		default:
		}
	}
	waitGroup.Done()
}


