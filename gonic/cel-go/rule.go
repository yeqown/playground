package main

type Rule interface {
	Id() string
	RuleDesc() string
}

type doctorRule struct {
	id, ruleDesc string
}

func mockDoctorRule(rid, ruleDesc string) Rule {
	return doctorRule{
		id:       rid,
		ruleDesc: ruleDesc,
	}
}

func (d doctorRule) Id() string {
	return d.id
}

func (d doctorRule) RuleDesc() string {
	return d.ruleDesc
}
