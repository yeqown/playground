package basic_test

import "testing"

// 测试目标：结构体中的内嵌结构体的修改
// case1: outterStruct1中的embed结构体的修改
// case2: outterStruct2中的embed结构体(指针)的修改
// case3: outterStruct1中的embed结构体, 在函数中中修改(传入指针)
// case4: outterStruct2中的embed结构体, 在函数中中修改（传入指针）

type embed struct {
	A string
}

type outterStruct1 struct {
	embed
	B string
}

type outterStruct2 struct {
	*embed
	B string
}

func modifyStruct1Embed(a *outterStruct1) {
	a.A = "hello2"
	a.B = "world2"
}

func modifyStruct2Embed(a *outterStruct2) {
	a.A = "hello2"
	a.B = "world2"
}

func Test_modify_struct1(t *testing.T) {
	var o1 = outterStruct1{
		embed: embed{A: "1111"},
		B:     "2222",
	}
	o1.A = "hello"
	o1.B = "world"
	if o1.A != "hello" || o1.B != "world" {
		t.Errorf("Error modify struct1")
	}

	// PASS
}

func Test_modify_struct2(t *testing.T) {
	var o2 = &outterStruct2{
		embed: &embed{},
		B:     "",
	}

	o2.A = "hello"
	o2.B = "world"
	if o2.A != "hello" || o2.B != "world" {
		t.Errorf("Error modify struct2")
	}

	// PASS
}

func Test_struct1_modifyStruct1Embed(t *testing.T) {
	var o1 = &outterStruct1{
		embed: embed{A: "1111"},
		B:     "2222",
	}

	modifyStruct1Embed(o1)

	if o1.A != "hello2" || o1.B != "world2" {
		t.Errorf("Error modify struct1")
	}

	// PASS
}

func Test_struc2_modifyStruct2Embed(t *testing.T) {
	var o2 = &outterStruct2{
		embed: &embed{},
		B:     "",
	}

	modifyStruct2Embed(o2)

	if o2.A != "hello2" || o2.B != "world2" {
		t.Errorf("Error modify struct2")
	}

	// PASS
}
