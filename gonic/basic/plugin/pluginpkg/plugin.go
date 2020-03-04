package main

import (
	"fmt"
)

var (
	// TestExportInt 用于测试plugin导出变量
	TestExportInt int
	// testUnexportedString 用于测试plugin导出小写变量
	testUnexportedString string
)

// 测试init函数在插件模式中是否可用
func init() {
	TestExportInt = 99
	testUnexportedString = "unexported string"
}

func main() {}

// CustomMethodName 用于测试函数导出
func CustomMethodName(name string) string {
	return fmt.Sprintf("this is exported name: %s", name)
}

// unexportedMethodName 用于测试函数导出
func unexportedMethodName(name string) string {
	return fmt.Sprintf("this is unexported name: %s", name)
}
