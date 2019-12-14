package services

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strings"
	"github.com/hpcloud/tail"
	"src/conf"
)

type TailWithConf struct {
	tail     *tail.Tail
	offset   int64
	logConf  conf.LogConfig
	secLimit *conf.SecondLimit
	exitChan chan bool
}

func (t *TailWithConf) readLog(fileName string) {

	for line := range t.tail.Lines {
		if line.Err != nil {
			logs.Error("read line error:%v ", line.Err)
			continue
		}

		lineStr := strings.TrimSpace(line.Text)
		fmt.Println("从被监控的文件", fileName, "中读到的字符串=", lineStr)

		if len(lineStr) == 0 || lineStr[0] == '\n' {
			continue
		}
		fmt.Println("向kafka生产者数据通道发送消息 消息字符串=",
			line.Text, "消息的topic=", t.logConf.Topic)
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

