package basic_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func Test_deferCall(t *testing.T) {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()
	panic("触发异常")

	/*
		output:
			打印后
			打印中
			打印前
			panic: 触发异常
	*/

	/*
		explain:
	*/

}

type student struct {
	Age  int
	Name string
}

func Test_parseStudent(t *testing.T) {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _, stu := range stus {
		m[stu.Name] = &stu
		// stu变量地址是固定的
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func Test_runtimeQuestion(t *testing.T) {
	// 这个还不清楚
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("j: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("showB")
}

func (p *People) showO() {
	fmt.Println("show0")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func Test_derived(t *testing.T) {
	t2 := Teacher{}
	t2.ShowA()
	t2.showO()

	// output:
	// showA
	// showB
	// show0

	// explain:
	// 类似于继承，调用
}

func Test_panicTest(t *testing.T) {
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "codegen"
	// 随机选择
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}
}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func Test_deferCalc(t *testing.T) {
	// defer 是把执行的函数计算好，然后放到先入后出的队列中
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	a = 1
	// 这里换成 a = 1 更有迷惑一点？
}

// 有语法错误吗？输出结果是啥。。。？
func Test_appendMake(t *testing.T) {
	s := make([]int, 5)
	arr := []int{1, 2, 3}
	s = append(s, arr...)
	fmt.Println(s)
}

type Human interface {
	Speak(string) string
}

type Student struct{}

func (stu *Student) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func Test_interfaceTest(t *testing.T) {
	var peo Student = Student{}
	think := "bitch"
	fmt.Println(peo.Speak(think))
}
