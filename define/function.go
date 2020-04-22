package define

type Consumer func(v interface{})

type Function func(t interface{}) (r interface{})

func testFunction() {

}

type F = Function

type I int64

type Int = int64
