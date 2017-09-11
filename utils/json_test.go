package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestJsonDecode2Map(t *testing.T) {
	json := `{"a":1,"location":{"0":{"province":"p1","city":"c1"},"1":[{"province":"p2","city":"c2"}],"range_list":[{"min":10,"max":100},{"max":1000}]},"range":{"min":0,"max":"9999"},"list":[1,2,3]}`
	res, err := JsonDecode2Map([]byte(json))
	fmt.Println(res, err)

	fmt.Println(reflect.TypeOf(res))

	for k, v := range res {
		if m, ok := v.(map[string]interface{}); ok {
			for k, v := range m {
				fmt.Println(k, reflect.TypeOf(v))
			}
		} else {
			fmt.Println(k, reflect.TypeOf(v))
		}
	}

	var val int64
	val = 10000
	fmt.Println(int(val))
}

func TestJsonEncode(t *testing.T) {
	data := map[string]interface{}{
		"name": "chen",
		"age": 18,
	}
	fmt.Println(JsonEncode(data))

}
