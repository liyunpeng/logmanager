package conf

import (
	"sync/atomic"
	"time"
	"github.com/astaxie/beego/logs"
)
// SecondLimit to Limit num in one second
type SecondLimit struct {
	unixSecond int64
	curCount   int32
	Limit      int32
}

// NewSecondLimit to init a SecondLimit obj
func NewSecondLimit(limit int32) *SecondLimit {
	secLimit := &SecondLimit{
		unixSecond: time.Now().Unix(),
		curCount:   0,
		Limit:      limit,
	}

	return secLimit
}

// Add is func to 
func (s *SecondLimit) Add(count int) {
	sec := time.Now().Unix()
	if sec == s.unixSecond {
		atomic.AddInt32(&s.curCount, int32(count))
		return
	}

	atomic.StoreInt64(&s.unixSecond, sec)
	atomic.StoreInt32(&s.curCount, int32(count))
}

// Wait to Limit num
func (s *SecondLimit) Wait() bool {
	for {
		sec := time.Now().Unix()
		if (sec == atomic.LoadInt64(&s.unixSecond)) && s.curCount >= s.Limit {
			time.Sleep(time.Millisecond)
			logs.Debug("Limit is runing, Limit: %d s.curCount:%d", s.Limit, s.curCount)
			continue
		}

		if sec != atomic.LoadInt64(&s.unixSecond) {
			atomic.StoreInt64(&s.unixSecond, sec)
			atomic.StoreInt32(&s.curCount, 0)
		}
		logs.Debug("Limit is exited")
		return false
	}
}
