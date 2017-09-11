package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	currentTime int64
	nowTime     time.Time
)

func init() {
	nowTime = time.Now()
	currentTime = nowTime.Unix()
	rand.Seed(currentTime)
	go tick()
}

const (
	SecondsOfHour = 3600
	SecondsOfDay  = 86400
	SecondsOfWeek = 604800
)

//使用频率最高, 采用弱精度时间
func UnixTime() int64 {
	return currentTime
}

//精确到秒
func NowTime() time.Time {
	return nowTime
}

func UnixTimeNano() int64 {
	return time.Now().UnixNano()
}

func UnixTimeMicro() int64 {
	return time.Now().UnixNano() / 1000
}

func tick() {
	for {
		select {
		case <-time.Tick(time.Second):
			currentTime = time.Now().Unix()
			nowTime = time.Now()
			rand.Seed(currentTime)
		}
	}
}

func FuncCost(f func()) int64 {
	start := time.Now().UnixNano()
	f()
	return (time.Now().UnixNano() - start) / 1000000
}

func TimerTask(task func() error, interval time.Duration) {
	defer RecoverFunc(task, true)
	for {
		select {
		case <-time.Tick(interval):
			if task == nil {
				break
			}
			if err := task(); err != nil {
				return
			}
		}
	}
}

func GetYear(timestamp int64) int {
	t := time.Unix(timestamp, 0)
	return t.Year()
}

func StartOfYear(year int) int {
	format := fmt.Sprintf("%d-01-01", year)
	if t, err := time.Parse("2006-02-01", format); err == nil {
		return int(t.Unix())
	}
	return 0
}

/*
%Y year
%m month
%d day
%H hour
%i minute
%s second
*/
func Str2Unix(fmt string) int64 {
	if t, err := time.Parse("2006-02-01", fmt); err == nil {
		return t.Unix()
	}
	return 0
}
