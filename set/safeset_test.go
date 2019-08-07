package set

import (
	"fmt"
	"testing"
)

func TestNewHashSet(t *testing.T) {
	var list []string
	list = append(list, "a")
	list = append(list, "b")
	list = append(list, "c")
	list = append(list, "d")
	list = append(list, "e")
	list = append(list, "f")

	set := NewSafeSet(list)

	fmt.Println(set)
	var res bool
	res = set.Has("a")
	fmt.Println(res)

	set.Remove("a")
	res = set.Has("a")
	fmt.Println(res)
}

func TestHashSet_Adds(t *testing.T) {
	var list []string
	list = append(list, "a")
	list = append(list, "b")
	list = append(list, "c")
	list = append(list, "d")
	list = append(list, "e")
	list = append(list, "f")

	set := NewSafeSet(nil)
	set.AddSlice(list)
	fmt.Println(set.Has("f"))
}
