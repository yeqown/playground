package basic_test

import (
	"testing"

	"github.com/casbin/casbin"
	mysqladpter "github.com/casbin/mysql-adapter"
)

var (
	casbinModelACL = `
# Request definition
[request_definition]
r = sub, obj, act

# Policy definition
[policy_definition]
p = sub, obj, act

# Policy effect
[policy_effect]
e = some(where (p.eft == allow))

# Matchers
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`
	e *casbin.Enforcer
)

func Test_CasbinACL(t *testing.T) {
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.

	a := mysqladpter.NewAdapter("mysql", "root:@tcp(127.0.0.1:3306)/")
	m := casbin.NewModel(casbinModelACL)
	e := casbin.NewEnforcer(m, a)

	// e.AddFunction(name, function)

	if ok := e.Enforce(sub, obj, act); ok {
		// pass to do something
	}
}
