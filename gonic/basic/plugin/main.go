package main

import (
	"fmt"
	"os"
	"plugin"
)

func main() {

	var (
		plg *plugin.Plugin
		err error

		soPath = "./pluginpkg/plugin.so"
	)
	if plg, err = plugin.Open(soPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// symbol is an interface in "plugin"
	// it can be varible or func

	// CustomMethodName is func
	if symFunc1, err := plg.Lookup("CustomMethodName"); err != nil {
		fmt.Printf("could not lookup CustomMethodName: %v\n", err)
	} else {
		s := symFunc1.(func(string) string)("CustomMethodName")
		fmt.Println(s)
	}

	// unexportedMethodName is func, could not found
	// if symFunc2, err := plg.Lookup("unexportedMethodName"); err != nil {
	// 	fmt.Printf("could not lookup unexportedMethodName: %v\n", err)
	// } else {
	// 	s := symFunc2.(func(string) string)("unexportedMethodName")
	// 	fmt.Println(s)
	// }

	// TestExportInt var int
	if var1, err := plg.Lookup("TestExportInt"); err != nil {
		fmt.Printf("could not lookup TestExportInt:: %v", err)
	} else {
		// why exported is (*int) as it described in plugin.go
		fmt.Println(*var1.(*int))
	}

	// testUnexportedString var string, could not found
	// if var2, err := plg.Lookup("testUnexportedString"); err != nil {
	// 	fmt.Printf("could not lookup testUnexportedString: %v", err)
	// } else {
	// 	fmt.Println(var2.(string))
	// }
}
