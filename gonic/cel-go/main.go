package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
)

func main() {
	d := &driver{
		ruleCache: map[string]cel.Program{},
	}

	obj := mockDoctorObj(12, 32, "med@y")
	rule1 := mockDoctorRule("exp1", "doctor.age > 32 || doctor.name.startsWith('med')")
	out, _ := d.Run(obj, rule1)
	fmt.Printf("exp1 executing: %v\n", out)
	// true

	rule2 := mockDoctorRule("exp2", "doctor.ordered(14) && doctor.id > uint(0)")
	_ = rule2
	out2, _ := d.Run(obj, rule2)
	fmt.Printf("exp2 executing: %v\n", out2)
	// false
}

type driver struct {
	ruleCache map[string]cel.Program
}

func (d *driver) Run(obj object, rule Rule) (interface{}, error) {
	// 构建表达式
	prg, ok := d.ruleCache[rule.Id()]
	if !ok {
		var err error
		prg, err = d.constructRule(obj, rule)
		if err != nil {
			println("could not constructRule: " + err.Error())
			return nil, err
		}
	}

	// 运行，表达式求值
	out, details, err := prg.Eval(map[string]interface{}{
		obj.key(): obj.proto(),
	})
	fmt.Printf("err=%v, out=%v, details=%v\n", err, out, details)
	return out.Value(), nil
}

func (d *driver) constructRule(obj object, rule Rule) (cel.Program, error) {
	envOptions := obj.options()
	env, err := cel.NewEnv(
		// come built in env
		envOptions...)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
		return nil, err
	}

	ast, iss := env.Compile(rule.RuleDesc())
	if iss.Err() != nil {
		log.Fatalln("check error:", iss.Err())
		return nil, err
	}

	//checked, _ := cel.AstToCheckedExpr(ast)
	//cel.AstToString()
	//parser.Unparse()

	fs := obj.functions()
	prg, err := env.Program(ast, cel.EvalOptions(cel.OptOptimize), cel.Functions(fs...))
	if err != nil {
		log.Fatalln("program error:", err)
		return nil, err
	}
	d.ruleCache[rule.Id()] = prg

	return prg, nil
}
