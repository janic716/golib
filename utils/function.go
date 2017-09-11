package utils

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"
	"time"
)

//执行指定函数f, 如果发生错误, 则重试
func RetryExecute(f func() error, retryTimes int) bool {
	return RetryExecuteWithWait(f, retryTimes, 0)
}

func RetryExecuteWithWait(f func() error, retryTimes int, msWaitTime int) bool {
	if retryTimes <= 0 {
		return false
	}
	var err error
	for i := 0; i < retryTimes; i++ {
		if err = f(); err == nil {
			return true
		}
		if msWaitTime > 0 {
			time.Sleep(time.Millisecond * time.Duration(msWaitTime))
		}
	}
	return false
}

func RetryExecuteWithTimeout(f func() error, retryTimes int, msTimeout int) bool {
	job := make(chan bool)
	go func() {
		job <- RetryExecuteWithWait(f, retryTimes, 0)
	}()
	select {
	case <-time.After(time.Millisecond * time.Duration(msTimeout)):
		return false
	case <-job:
		return true
	}
}

func Execute(f func() error, times int) int {
	var (
		err error
		res int
	)
	for i := 0; i < times; i++ {
		if err = f(); err == nil {
			res++
		}
	}
	return res
}

func RecoverFunc(p interface{}, logStack bool) {
	if r := recover(); r != nil {
		t := reflect.TypeOf(p)
		if logStack {
			fmt.Printf("%s panic, stack: %s", t.String(), string(debug.Stack()))
		} else {
			fmt.Printf("%s panic", t.String())
		}
	}
}

//map can execute function f in the given num of concurrency
//list must be slice, any item in it  is the param of f
func Map(f func(interface{}) interface{}, list []interface{}, maxConcurrency int) []interface{} {
	var waitGroup sync.WaitGroup
	jobs := make(chan int, maxConcurrency)
	started := make(chan bool)
	resList := make([]interface{}, len(list), len(list))
	for i, v := range list {
		resList[i] = nil
		waitGroup.Add(1)
		jobs <- i
		index := i
		tmp := v
		go func(int, interface{}) {
			started <- true
			defer func() {
				waitGroup.Done()
				<-jobs
			}()
			resList[index] = f(v)
		}(index, tmp)
		<-started
	}
	waitGroup.Wait()
	return resList
}

func Reduce(f func(interface{}, interface{}) interface{}, list []interface{}) interface{} {
	var res interface{}
	if list == nil || len(list) == 0 {
		return res
	}
	for _, val := range list {
		if val == nil {
			continue
		}
		if res == nil {
			res = val
			continue
		}
		res = f(res, val)
	}
	return res
}

func ReduceConcurrency(f func(interface{}, interface{}) interface{}, list []interface{}, maxConcurrency int) interface{} {
	listCount := len(list)
	if listCount == 0 {
		return nil
	}
	if listCount == 1 {
		return list[0]
	}
	var newList []interface{}
	var waitGroup sync.WaitGroup
	var mutex sync.Mutex
	jobs := make(chan int, maxConcurrency)
	started := make(chan bool)
	var first interface{}
	var second interface{}
	for i := 0; i < len(list); i++ {
		if list[i] == nil {
			continue
		}
		if first == nil {
			first = list[i]
			continue
		}
		if second == nil {
			second = list[i]
		}
		jobs <- i
		waitGroup.Add(1)
		v1, v2 := first, second
		go func(val1 interface{}, val2 interface{}) {
			started <- true
			defer func() {
				<-jobs
				waitGroup.Done()
			}()
			if val1 == nil || val2 == nil {
				//fmt.Println(val1, val2)
			}
			r := f(val1, val2)
			if r != nil {
				mutex.Lock()
				newList = append(newList, r)
				mutex.Unlock()
			}

		}(v1, v2)
		<-started
		first, second = nil, nil
	}
	waitGroup.Wait()
	if first != nil {
		newList = append(newList, first)
	}
	return ReduceConcurrency(f, newList, maxConcurrency)
}

func ReduceSliceMergeFunc(val1 interface{}, val2 interface{}) interface{} {
	if val1 == nil && val2 != nil {
		return val2
	}
	if val2 == nil && val1 != nil {
		return val1
	}
	tp1 := reflect.TypeOf(val1)
	tp2 := reflect.TypeOf(val2)
	var res []interface{}

	if tp1.Kind() == reflect.Slice && tp1.Kind() == reflect.Slice {
		value1 := reflect.ValueOf(val1)
		value2 := reflect.ValueOf(val2)
		if value1.Len() > 0 {
			v1 := ToInterfaceList(value1.Interface())
			fmt.Println("v1", v1)
			res = append(res, ToInterfaceList(value1.Interface())...)
		}
		if value2.Len() > 0 {
			v2 := ToInterfaceList(value2.Interface())
			fmt.Println("v2", v2)
			res = append(res, ToInterfaceList(value2.Interface())...)
		}
	} else {
		panic(fmt.Sprintf("type is err, val1:%s, %v, val2:%s, %v", tp1, val1, tp2, val2))
	}
	return res
}

func ReduceSliceIntMergeFunc(val1, val2 interface{}) interface{} {
	if val1 == nil {
		return val2
	}
	if val2 == nil {
		return val1
	}
	if v1, ok := val1.([]int); ok {
		if v2, ok := val2.([]int); ok {
			return append(v1, v2...)
		}
	}
	panic(fmt.Sprintf("type is err, val1:%s, %v, val2:%s, %v", reflect.TypeOf(val1), val1, reflect.TypeOf(val2), val2))

}
