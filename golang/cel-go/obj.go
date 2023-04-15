package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	pbexpr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	proto2 "google.golang.org/protobuf/proto"
)

type object interface {
	key() string
	proto() proto2.Message
	options() []cel.EnvOption
	functions() []*functions.Overload
}

func mockDoctorObj(id, age int, name string) object {
	return &Doctor{
		Id:   uint64(id),
		Name: name,
		Age:  int64(age),
	}
}

func (m *Doctor) proto() proto2.Message {
	return proto.MessageV2(m)
}

func (m *Doctor) key() string {
	return "doctor"
}

func (m *Doctor) options() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Types(proto.MessageV2(m)), // 注册类型
		cel.Declarations(
			decls.NewVar("doctor", decls.NewObjectType("celgo.demo.Doctor")),
			decls.NewFunction("ordered", decls.NewParameterizedInstanceOverload(
				"doctor_ordered",
				//[]*pbexpr.Type{decls.NewTypeParamType("doctorId"), decls.NewTypeParamType("n")},
				[]*pbexpr.Type{decls.NewObjectType("celgo.demo.Doctor"), decls.NewTypeParamType("n")},
				decls.Bool,
				[]string{"doctorId", "n"},
			)),
		), // 声明变量和函数
	}
}

func (m *Doctor) functions() []*functions.Overload {
	return []*functions.Overload{
		{
			Operator: "doctor_ordered",
			Binary:   m.genOrdered(),
		},
	}
}

// doctor.ordered(10) || other conditions
func (m *Doctor) genOrdered() functions.BinaryOp {
	return func(lhs ref.Val, rhs ref.Val) ref.Val {
		i, ok := rhs.(types.Int)
		if !ok {
			return types.NewErr("no such overload: not int")
		}

		// println("called: doctorId=", doctorId, "n=",n)
		if m.Id > uint64(i) {
			return types.True
		}

		return types.False
	}
}
