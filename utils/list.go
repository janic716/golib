package utils

import (
	"reflect"
	"strings"
)

//func ListChunk(list interface{}, size int, pchunks interface{}) error {
//	listValue := reflect.Indirect(reflect.ValueOf(list))
//	if listValue.Kind() != reflect.Slice {
//		return errors.New(fmt.Sprintf("%v is not slice type", list))
//	}
//	tp := reflect.TypeOf(pchunks)
//	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Slice {
//		return errors.New("list must be slice ptr")
//	}
//	itemType := tp.Elem().Elem()
//	chunksValue := reflect.ValueOf(pchunks).Elem()
//	count := listValue.Len()
//	for i:=0; i < count; i++ {
//		value := reflect.New(tp).Elem()
//		valueList = reflect.Append(valueList, value)
//	}
//}

func Implode(list interface{}, seq string) string {
	listValue := reflect.Indirect(reflect.ValueOf(list))
	if listValue.Kind() != reflect.Slice {
		return ""
	}
	count := listValue.Len()
	listStr := make([]string, 0, count)
	for i := 0; i < count; i++ {
		v := listValue.Index(i)
		if str, err := GetValue(v); err == nil {
			listStr = append(listStr, str)
		}
	}
	return strings.Join(listStr, seq)
}

func ListChunks(list interface{}, size int) []interface{} {
	var res []interface{}
	value := reflect.Indirect(reflect.ValueOf(list))
	if value.Kind() == reflect.Slice {
		count := value.Len()
		for i := 0; i < count; i += size {
			res = append(res, value.Slice(i, MinInt(i+size, count)).Interface())
		}
	}
	return res
}

func ToInterfaceList(list interface{}) []interface{} {
	var res []interface{}
	value := reflect.Indirect(reflect.ValueOf(list))
	if value.Kind() == reflect.Slice {
		count := value.Len()
		for i := 0; i < count; i++ {
			res = append(res, value.Index(i).Interface())
		}
	}
	return res
}
