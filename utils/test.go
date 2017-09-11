package utils

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

func TestFunc(f func(), times int) (info string) {
	start := time.Now().UnixNano()
	for i := 0; i < times; i++ {
		f()
	}
	end := time.Now().UnixNano()
	totalTime := float64(end-start) / 1000
	avg := float64(totalTime) / float64(times)
	unit := "us"
	if avg > 100000 {
		totalTime = totalTime / 1000
		avg = avg / 10000
		unit = "ms"
	}
	res := fmt.Sprintf("execute-times:%d | total-time:%.3f %s | avg-time:%.3f %s ", times, totalTime, unit, float64(totalTime)/float64(times), unit)
	return res
}

func AssertMust(b bool) {
	if b {
		return
	}
	panic("assertion failed")
}

func AssertMustString(real, expect string) {
	panicIfNoTrue(real == expect, real, expect)
}

func AssertEqual(real, expect interface{}) {
	realStr, expectStr := ValueToString(real), ValueToString(expect)
	if realStr != expectStr {
		panic(fmt.Sprintf("assertion failed, real: %v, expect: %v", real, expect))
	}
}

func AssertStrictEqual(real, expect interface{}) {
	if reflect.TypeOf(real).Kind() != reflect.TypeOf(expect).Kind() {
		panic(fmt.Sprintf("assertion failed, %v %v type is not equal", real, expect))
	}
	realStr, expectStr := ValueToString(real), ValueToString(expect)
	if realStr != expectStr {
		panic(fmt.Sprintf("assertion failed, real: %v, expect: %v", real, expect))
	}
}

func ValueToString(val interface{}) string {
	tp := reflect.TypeOf(val)
	switch tp.Kind() {
	case reflect.String, reflect.Uint, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
		return fmt.Sprint(val)
	case reflect.Map:
		value := reflect.ValueOf(val)
		keys := value.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return strings.Compare(keys[i].String(), keys[j].String()) > 0
		})
		var list []string
		for _, k := range keys {
			v := value.MapIndex(k)
			list = append(list, v.String())
		}
		return strings.Join(list, "|")
	default:
		return fmt.Sprint(val)
	}
	return ""
}

func AssertMustInt(real, expect int) {
	panicIfNoTrue(real == expect, real, expect)
}

func panicIfNoTrue(b bool, real, expect interface{}) {
	if !b {
		panic(fmt.Sprintf("assertion failed, real: %v, expect: %v", real, expect))
	}
}

func AssertMustNoError(err error) {
	if err == nil {
		return
	}
	panic(fmt.Sprintf("%s error happens, assertion failed", err.Error()))
}

func DeepEqual(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}
