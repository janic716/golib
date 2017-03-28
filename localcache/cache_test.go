package localcache

import (
	"testing"
	"time"
	"fmt"
	"strconv"
	"math/rand"
)

func TestGetInt(t *testing.T) {
	var (
		val int
		err error
	)
	key := "abc"
	Set(key, -100, 3)
	if val, err := GetInt(key); err == nil {
		if val != -100 {
			panic("get err")
		}
	} else {
		panic(err)
	}
	Del(key)
	if _, err := GetInt(key); err == nil {
		panic("del err")
	} else {
		fmt.Println(err)
	}

	Set(key, -100, 3)
	val, err = GetInt(key)
	fmt.Println(val, err)
	time.Sleep(time.Duration(time.Second * 4))
	fmt.Println(time.Now().Unix())
	val, err = GetInt(key)
	fmt.Println(val, err)
	if val, err := GetInt(key); err == nil {
		fmt.Println(val)
		panic("")
	}

	Set(key, 10000, 1)

	Set(key, "123", 2)
	if val, err := GetInt(key); err == nil {
		fmt.Println(val)
	} else {
		panic(err)
	}

	fmt.Println("GetInt pass")
}

func TestConcurrent(t *testing.T) {
	jobs := make(chan bool, 100)
	start := time.Now().Unix()
	jobCnt := 300
	fmt.Println(time.Now().UnixNano())

	for i := 0; i < jobCnt / 3; i++ {
		go func() {
			for j := 0; j < 200000; j++ {
				Set(strconv.FormatInt(time.Now().UnixNano() / 1000, 10), strconv.Itoa(int(rand.Int31())), 3)
				//time.Sleep(time.Duration(time.Nanosecond * 100000))
			}
			jobs <- true
		}()
	}
	for i := 0; i < jobCnt / 3; i++ {
		go func() {
			for j := 0; j < 200000; j++ {
				Get(strconv.FormatInt(time.Now().UnixNano() / 1000, 10))
				//time.Sleep(time.Duration(time.Nanosecond * 100000))

			}
			jobs <- true
		}()
	}
	for i := 0; i < jobCnt / 3; i++ {
		go func() {
			for j := 0; j < 200000; j++ {
				Del(strconv.FormatInt(time.Now().UnixNano() / 1000, 10))
				//time.Sleep(time.Duration(time.Nanosecond * 100000))

			}
			jobs <- true
		}()
	}
	for i := 0; i < jobCnt; i++ {
		<-jobs
	}
	fmt.Println(Stat())
	fmt.Println(time.Now().Unix() - start)
	ResetStat()
	fmt.Println(Stat())
	close(jobs)
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
	//time.Sleep(time.Duration(time.Second * 10))
	//fmt.Println(Stat())
}

func TestGetStrings(t *testing.T) {
	println(cache.get)
	res := []string{"a", "b", "c", "d", "e", "f"}
	Set("abc", res, 1000)
	r, err := GetStrings("abc")
	fmt.Println("r -> ", r, err)

	r2, found := Get("abc")
	fmt.Println(r2, found)
}

func TestGetFloat64(t *testing.T) {
	key := "float64"
	Set(key, 12131.131131, 10)
	if val, err := GetFloat64(key); err == nil {
		fmt.Println(val)
	}

	Set(key, 12131, 10)

	if val, err := GetFloat64(key); err == nil {
		fmt.Println(val)
	} else {
		fmt.Println("err", val)
	}

	in := "012131"
	Set(key, in, 10)
	if val, err := GetFloat64(key); err == nil {
		fmt.Println(val)
	} else {
		fmt.Println("err", val)
	}
}
