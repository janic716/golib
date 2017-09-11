package utils

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Name     string `orm:"name"`
	Age      int    `orm:"age"`
	Location `orm:"location"`
}

type PUser struct {
	Name *string
	Age  *int
	PLocation
}

type PLocation struct {
	Province *string `orm:"province"`
	City     *string `orm:"city"`
	Name     *string `orm:"user"`
}

type Location struct {
	Province string `orm:"province"`
	City     string `orm:"city"`
	Name     string `orm:"user"`
}

type Manager struct {
	User
	Title    string `orm:"my_title"`
	Level    int32
	ForgetMe string `orm:"-"`
	salary   int
}

func TestType(t *testing.T) {
	usr1 := User{}
	usr2 := User{}

	fmt.Println(reflect.TypeOf(usr1) == reflect.TypeOf(usr2))
}

func TestValueToString(t *testing.T) {
	data := make(map[string]string)
	data["name"] = "chenjian"
	data["age"] = fmt.Sprint(20)
	data["province"] = "handong"
	data["city"] = "jingzhou"
	data["my_title"] = "super rd"
	data["level"] = fmt.Sprint(2)
	var res string
	res = ValueToString(data)
	fmt.Println("res:", res)

	res = ValueToString(data)
	fmt.Println("res:", res)

	res = ValueToString(data)
	fmt.Println("res:", res)
}

func TestGetStructTagMap(t *testing.T) {
	if m, err := GetStructTagMap("orm", reflect.TypeOf(&Manager{})); err == nil {
		fmt.Println(m)
		AssertMustString(m["City"], "city")
		AssertMustString(m["Age"], "age")
		AssertMustString(m["Province"], "province")
		AssertMustString(m["Name"], "name")
		AssertMustString(m["Title"], "my_title")

	} else {
		panic(err)
	}
}

func TestToStruct(t *testing.T) {
	data := make(map[string]string)
	data["name"] = "chenjian"
	data["age"] = "28"
	data["province"] = "handong"
	data["city"] = "jingzhou"
	data["my_title"] = "super rd"
	data["level"] = "2"
	data["title"] = "great rd"
	var manage Manager
	err := ToStruct(data, &manage, "orm")
	fmt.Println(manage)
	if err != nil {
		panic(err)
	}
	res := ListStructFiles(manage)
	fmt.Println(res)
}

func TestToStructList(t *testing.T) {
	dataList := make([]map[string]string, 0)
	num := 10
	for i := 0; i < num; i++ {
		data := make(map[string]string)
		data["name"] = "chenjian"
		data["age"] = fmt.Sprint(i + 20)
		data["province"] = "handong"
		data["city"] = "jingzhou"
		data["my_title"] = "super rd"
		data["level"] = fmt.Sprint(i)
		dataList = append(dataList, data)
	}

	var list []Manager
	err := ToStructList(dataList, &list, "orm")
	if err != nil {
		panic(err)
	}
	AssertMustInt(len(list), num)
	for i, v := range list {
		AssertEqual(v.Name, "chenjian")
		AssertEqual(v.Age, i+20)
		AssertEqual(v.Province, "handong")
		AssertEqual(v.City, "jingzhou")
		AssertEqual(v.Title, "super rd")
		AssertEqual(v.Level, i)

	}
}

func TestSetValue(t *testing.T) {
	var list []interface{}
	list = append(list, "a")

	for _, v := range list {
		value := reflect.ValueOf(v)
		SetValue(&value, "abc")
		fmt.Println(v)
	}
	fmt.Println(list)
}

func TestToMap(t *testing.T) {
	var manager Manager
	manager.Name = "chenjian"
	manager.Age = 100
	manager.Title = "hello world"
	manager.Province = "province"
	manager.City = "city"
	manager.Level = 10
	manager.ForgetMe = "我会被忽略"

	data, err := ToMap(manager, "orm")
	fmt.Println(data, err)

	data, err = ToMap(&manager, "orm")
	fmt.Println(data, err)
}

func TestToStruct2(t *testing.T) {
	data := make(map[string]string)
	data["name"] = "chenjian"
	data["age"] = "28"
	data["province"] = "handong"
	data["city"] = "jingzhou"
	data["my_title"] = "super rd"
	data["level"] = "2"
	data["title"] = "great rd"
	var puser PUser
	err := ToStruct(data, &puser, "orm")
	fmt.Println(*puser.Age, *puser.Name, *puser.Province, *puser.City)
	if err != nil {
		panic(err)
	}

}

func TestToMap2(t *testing.T) {
	var puser PUser
	name := "chenjian"
	province := "province"
	city := "city"
	age := 100

	puser.Name = &name
	puser.Age = &age
	puser.Province = &province
	puser.City = &city

	data, err := ToMap(puser, "orm")
	fmt.Println(data, err)

	//data, err = ToMap(&puser, "orm")
	//fmt.Println(data, err)
}

func TestGetValue(t *testing.T) {
	var str string
	var vint int
	var pstr *string
	str = "hello"
	pstr = &str
	vint = 100
	fmt.Println(GetValue(reflect.ValueOf(str)))
	fmt.Println(GetValue(reflect.ValueOf(vint)))
	fmt.Println(GetValue(reflect.ValueOf(pstr)))
	user := User{Name: "chenjian"}
	fmt.Println(GetValue(reflect.ValueOf(user)))
}

func TestGetValueByFieldName(t *testing.T) {
	var manager Manager
	manager.Name = "chenjian"
	manager.Age = 100
	manager.Title = "hello world"
	manager.Province = "province"
	manager.City = "city"
	manager.Level = 10
	manager.ForgetMe = "我会被忽略"
	value := reflect.ValueOf(manager)
	fmt.Println(GetStructValueByName(manager, "Level"))
	fmt.Println(GetStructValueByName(manager, "ForgetMe"))
	fmt.Println(GetStructValueByName(&manager, "Name"))
	fmt.Println(GetStructValueByName(&manager, "salary"))
	fmt.Println(GetStructValueByName(&manager, "Salary"))
	fmt.Println(value.FieldByName("Name").Interface())
	var res string

	res = TestFunc(func() {
		value.FieldByName("Level").Interface()
		value.FieldByName("ForgetMe").Interface()
		value.FieldByName("Name").Interface()

	}, 10000)
	fmt.Println(res)

	res = TestFunc(func() {
		GetStructValueByName(value, "Level")
		GetStructValueByName(value, "ForgetMe")
		GetStructValueByName(value, "Name")
	}, 10000)
	fmt.Println(res)

	if v, ok := GetStructValueByName(manager, "Level").(int); ok {
		fmt.Println(v)
	}

	//res = utils.TestFunc(func() {
	//	fmt.Println(manager.Level)
	//	fmt.Println(manager.ForgetMe)
	//	fmt.Println(manager.Name)
	//
	//}, 10)
	//fmt.Println(res)
}
