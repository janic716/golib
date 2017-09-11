package utils

import (
	"fmt"
	"testing"
)

func TestStr2Unix(t *testing.T) {
	var res int64
	res = Str2Unix("1900-01-01")
	fmt.Println(res)
}
