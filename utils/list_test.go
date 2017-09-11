package utils

import (
	"fmt"
	"testing"
)

func TestImplode(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7}
	res := Implode(list, ",")
	fmt.Println(res)

	list1 := []int16{1, 2, 3, 4, 5, 6, 7}
	res = Implode(list1, ",")
	fmt.Println(res)

	list2 := []float64{1.5, 2.1, 3.0, 4, 5.9, 6.7, 7.7}
	res = Implode(list2, ",")
	fmt.Println(res)

	var list3 []float64
	res = Implode(list3, ",")
	fmt.Println("res: ", res)

	list4 := make([]interface{}, 4, 4)
	list4[0] = "str"
	list4[1] = 19
	list4[2] = 19.999
	list4[3] = true
	res = Implode(list4, ",")
	fmt.Println(res)

}

func TestListChunks(t *testing.T) {
	var list []int
	for i := 0; i < 23; i++ {
		list = append(list, i)
	}
	res := ListChunks(list, 10)
	fmt.Println(res)
}
