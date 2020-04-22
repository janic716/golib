package define

import "testing"

func Test_testFunction(t *testing.T) {
	f1 := func(r string) int64 {
		return 1
	}

	f2 := func(r string) interface{} {
		return 1
	}

	f3 := func(r interface{}) interface{} {
		return 1
	}

	var f Function
	f = f3
	println(f1, f2, f3)
	println(f)

	f4 := func(int64) {}

	var i I
	i = 64
	f4(int64(i))
	var j Int
	f4(j)

}
