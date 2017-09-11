package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	tagMapCache = make(map[string]map[string]string)
	mutex       sync.RWMutex
)

func ListStructFiles(stru interface{}) []interface{} {
	var res []interface{}
	value := reflect.Indirect(reflect.ValueOf(stru))
	kind := value.Kind()
	if kind == reflect.Invalid {
		return nil
	}
	tp := value.Type()
	count := tp.NumField()
	for i := 0; i < count; i++ {
		v := value.Field(i)
		if v.CanInterface() {
			res = append(res, v.Interface())
		}
	}
	return res
}

func GetStructTagMap(tag string, tp reflect.Type) (map[string]string, error) {
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("%v must be struct type", tp))
	}
	mutex.RLock()
	if v, found := tagMapCache[tp.String()]; found {
		mutex.RUnlock()
		return v, nil
	}
	mutex.RUnlock()
	tagMap := make(map[string]string)
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if f.Anonymous {
			if fMap, err := GetStructTagMap(tag, f.Type); err == nil {
				for k, v := range fMap {
					if _, found := tagMap[k]; !found {
						tagMap[k] = v
					}
				}
			}
		} else {
			tagVal := f.Tag.Get(tag)
			if tagVal == "" {
				tagMap[f.Name] = strings.ToLower(f.Name)
			} else {
				tagMap[f.Name] = tagVal
			}
		}
	}
	if tagMap != nil {
		mutex.Lock()
		tagMapCache[tp.String()] = tagMap
		mutex.Unlock()
	}
	return tagMap, nil
}

func ToStruct(data map[string]string, stru interface{}, ormTag string) error {
	tp := reflect.TypeOf(stru)
	value := reflect.ValueOf(stru).Elem()
	if tp.Kind() != reflect.Ptr {
		return errors.New("list must be ptr ptr")
	}
	tagMap, err := GetStructTagMap(ormTag, tp)
	if err != nil {
		return err
	}
	for k, v := range tagMap {
		fieldValue := value.FieldByName(k)
		err := SetValue(&fieldValue, data[v])
		if err != nil {
			return err
		}
	}
	return nil
}

func ToStruct2(data map[string]interface{}, stru interface{}, ormTag string) error {
	tp := reflect.TypeOf(stru)
	value := reflect.ValueOf(stru).Elem()
	if tp.Kind() != reflect.Ptr {
		return errors.New("list must be ptr ptr")
	}
	tagMap, err := GetStructTagMap(ormTag, tp)
	if err != nil {
		return err
	}
	for k, v := range tagMap {
		fieldValue := value.FieldByName(k)
		if val, found := data[v]; found {
			if err := SetValue(&fieldValue, ToString(val)); err != nil {
				return err
			}
		}
	}
	return nil
}

func ToString(val interface{}) string {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(float64(val.(float32)), 'f', -1, 32)
	default:
		return fmt.Sprint(val)
	}
	return ""

}

func ToMap(stru interface{}, ormTag string) (map[string]interface{}, error) {
	tp := reflect.TypeOf(stru)
	value := reflect.ValueOf(stru)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		value = value.Elem()
	}
	tagMap, err := GetStructTagMap(ormTag, tp)
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]interface{})
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		f := tp.Field(i)
		vf := value.Field(i)
		if !vf.CanInterface() {
			continue
		}
		if f.Anonymous {
			m, err := ToMap(vf.Interface(), ormTag)
			if err != nil {
				return nil, err
			}
			for k, v := range m {
				if _, found := dataMap[k]; !found {
					dataMap[k] = v
				}
			}
			continue
		}
		if v, found := tagMap[f.Name]; found {
			if v != "-" {
				if vf.Type().Kind() == reflect.Ptr {
					if vf.Elem().Kind() != reflect.Invalid {
						dataMap[v] = vf.Elem().Interface()
					}
				} else {
					dataMap[v] = vf.Interface()
				}
			}
		} else {
			dataMap[f.Name] = value.Field(i).Interface()
		}
	}
	return dataMap, nil
}

func GetValue(value reflect.Value) (res string, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		res, err = GetValue(value.Elem())
	default:
		res = fmt.Sprint(value.Interface())
	}
	return
}

func GetStructValueByName(stru interface{}, fieldName string) interface{} {
	value := reflect.ValueOf(stru)
	kind := value.Kind()
	if kind == reflect.Invalid {
		return nil
	}
	if kind == reflect.Ptr {
		value = value.Elem()
		return GetStructValueByName(value.Interface(), fieldName)
	}
	if f, ok := value.Type().FieldByName(fieldName); ok {
		v := value.FieldByIndex(f.Index)
		if v.CanInterface() {
			return v.Interface()
		}
	}
	return nil
}

func SetValue(value *reflect.Value, set string) error {
	if value.CanSet() && len(set) > 0 {
		switch value.Kind() {
		case reflect.String:
			value.SetString(set)
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			val, err := strconv.ParseInt(set, 10, 64)
			if err != nil {
				return err
			}
			value.SetInt(val)
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint8, reflect.Uint16:
			val, err := strconv.ParseUint(set, 10, 64)
			if err != nil {
				return err
			}
			value.SetUint(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(set)
			if err != nil {
				return err
			}
			value.SetBool(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(set, 64)
			if err != nil {
				return err
			}
			value.SetFloat(val)
		case reflect.Ptr:
			tp := value.Type().Elem()
			newValue := reflect.New(tp)
			value.Set(newValue)
			tmp := value.Elem()
			err := SetValue(&tmp, set)
			if err == nil {
				value.Elem().Set(tmp)
			}
		}
	}
	return nil
}

func ToStructList(dataList []map[string]string, list interface{}, ormTag string) error {
	if len(dataList) == 0 {
		return nil
	}
	tp := reflect.TypeOf(list)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Slice {
		return errors.New("list must be slice ptr")
	}
	tp = tp.Elem().Elem()
	valueList := reflect.ValueOf(list).Elem()
	tagMap, err := GetStructTagMap(ormTag, tp)
	if err != nil {
		return err
	}
	for _, data := range dataList {
		value := reflect.New(tp).Elem()
		for k, v := range tagMap {
			fieldValue := value.FieldByName(k)
			if _, found := data[v]; found {
				if err := SetValue(&fieldValue, data[v]); err != nil {
					return err
				}
			}
		}
		valueList = reflect.Append(valueList, value)
	}
	reflect.ValueOf(list).Elem().Set(valueList)
	return nil
}
