package utils

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	var list []interface{}
	for i := 0; i < 20000; i++ {
		list = append(list, i)
	}
	even := func(input interface{}) interface{} {
		if v, ok := input.(int); ok {
			//fmt.Println(v % 2)
			if v%2 == 0 {
				time.Sleep(time.Microsecond * 1000)
				return v
			}
		}
		return nil
	}
	var res string
	res = TestFunc(func() {
		Map(even, list, 20)
	}, 1)
	fmt.Println(res)
	res = TestFunc(func() {
		Map(even, list, 10)
	}, 1)
	fmt.Println(res)
	res = TestFunc(func() {
		Map(even, list, 2)
	}, 1)
	//Map(even, list, 10)
	//fmt.Println(list, res)
	fmt.Println(res)
}

func TestReduceConcurrency(t *testing.T) {
	var list []interface{}
	for i := 0; i < 1000001; i++ {
		list = append(list, i)
	}
	sum := func(val1 interface{}, val2 interface{}) interface{} {
		time.Sleep(time.Microsecond * 1)
		var v1, v2 int
		if v, ok := val1.(int); ok {
			v1 = v
		}
		if v, ok := val2.(int); ok {
			v2 = v
		}
		return v1 + v2
	}
	var res string
	res = TestFunc(func() {
		fmt.Println(ReduceConcurrency(sum, list, 2))

	}, 1)
	fmt.Println(res)
	res = TestFunc(func() {
		fmt.Println(ReduceConcurrency(sum, list, 10))

	}, 1)
	fmt.Println(res)
	res = TestFunc(func() {
		fmt.Println(ReduceConcurrency(sum, list, 20))

	}, 1)
	fmt.Println(res)
}

func TestReduce(t *testing.T) {
	var list []interface{}
	for i := 0; i < 100001; i++ {
		list = append(list, i)
	}
	list = append(list, "1")
	sum := func(val1 interface{}, val2 interface{}) interface{} {
		//time.Sleep(time.Microsecond * 10)
		var v1, v2 int
		if v, ok := val1.(int); ok {
			v1 = v
		}
		if v, ok := val2.(int); ok {
			v2 = v
		}
		return v1 + v2
	}
	var res string
	res = TestFunc(func() {
		fmt.Println(Reduce(sum, list))

	}, 1)
	fmt.Println(res)
}

func TestReduceSliceMerge(t *testing.T) {
	var list []interface{}
	for i := 0; i < 40; i++ {
		list = append(list, i)
	}
	list1 := ListChunks(list, 10)
	fmt.Println("len list1", len(list1))
	fmt.Println(list)
	doubble := func(input interface{}) interface{} {
		if v, ok := input.(int); ok {
			return v * 2
		}
		fmt.Println("nil", reflect.TypeOf(input))
		return nil
	}
	res := Map(func(val interface{}) interface{} {
		l := ToInterfaceList(val)
		res := Map(doubble, l, 1)
		fmt.Println("res", res)
		return res
	}, list1, 1)
	fmt.Println(res)
	res2 := Reduce(ReduceSliceMergeFunc, res)
	fmt.Println(res2)

	var listTest [][]int
	var listTest2 []interface{}
	for i := 0; i < 4; i++ {
		l := make([]int, 0)
		for j := 0; j < 10; j++ {
			l = append(l, (i*10+j)*2)
		}
		listTest = append(listTest, l)
		listTest2 = append(listTest2, l)
	}

	r1 := MergeSliceInt(listTest)
	fmt.Println(r1)
	r2 := Reduce(ReduceSliceMergeFunc, res)
	fmt.Println(r2)
	times := 100000
	var testRes string
	testRes = TestFunc(func() {
		r1 = MergeSliceInt(listTest)
	}, times)
	fmt.Println(testRes)

	testRes = TestFunc(func() {
		r2 = Reduce(ReduceSliceMergeFunc, res)
	}, times)
	fmt.Println(testRes)

	testRes = TestFunc(func() {
		r2 = Reduce(reduceMergeList, listTest2)
	}, times)
	fmt.Println(testRes)
}

func TestReduceSliceMergeFunc(t *testing.T) {
	//[[7 1000009] [] []]
	listData := make([]interface{}, 0, 100)
	listData = append(listData, []int{7, 1000009})
	listData = append(listData, []int{})
	listData = append(listData, []int{})
	listData = append(listData, []int{})

	fmt.Println(listData)
	res := Reduce(ReduceSliceIntMergeFunc, listData)
	fmt.Println(res, reflect.TypeOf(res))

}

func MergeSliceInt(val [][]int) []int {
	var res []int
	for _, v := range val {
		res = append(res, v...)
	}
	return res
}

func reduceMergeList(val1 interface{}, val2 interface{}) interface{} {
	if val1 == nil {
		return val2
	}
	if val2 == nil {
		return val1
	}
	if v1, ok := val1.([]int); ok {
		if v2, ok := val2.([]int); ok {
			v1 = append(v1, v2...)
			return v1
		}
	}
	return nil
}
