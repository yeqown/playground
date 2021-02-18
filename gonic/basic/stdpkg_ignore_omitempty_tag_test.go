package basic_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_ignore_omitempty(t *testing.T) {
	var a = &Page{
		A: 0,
		B: "",
	}

	v := omitemptyTag(a)
	dat, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("new %s\n", dat)

	dat2, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Printf("old: %s\n", dat2)
}

type Page struct {
	_ string `json:"ok"`
	A int    `json:"a,omitempty"`
	B string `json:"b,omitempty"`
}

func omitemptyTag(v interface{}) interface{} {
	ele := reflect.TypeOf(v).Elem()
	fields := []reflect.StructField{}

	for idx := 0; idx < ele.NumField(); idx++ {
		tag := ele.Field(idx).Tag.Get("json")
		fields = append(fields, ele.Field(idx))
		println("tag", tag)
		fields[idx].Tag = reflect.StructTag(strings.Replace(tag, ",omitempty", "", 1))
	}

	out := reflect.StructOf(fields)

	value := reflect.ValueOf(v).Elem()
	v2 := value.Convert(out)
	fmt.Println(v2.Interface(), v2.Type())

	return v2.Interface()
}
