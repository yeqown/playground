package tpl

import (
	"bytes"
	"text/template"
)

var (
	scaleResultTpl = `
{{ $d1 := index .scores 0 }}
{{ $d2 := index .scores 1 }}
{{ define "result1" }} 多动症 {{ end }}
{{ define "result2" }} 密集恐惧症 {{ end }}
{{ define "result3" }} 幽闭恐惧症 {{ end }}
{{if ge $d1 20}}得分：{{$d1}},{{$d2}}, {{- template "result1" }}
{{ else if ge $d2 20 }}得分：{{$d1}},{{$d2}}, {{- template "result2" }}
{{ else }}得分：{{$d1}},{{$d2}}, {{- template "result3" }}
{{ end }}
`
	tpl *template.Template
)

func init() {
	tpl = template.Must(
		template.New("tpl_example").Parse(scaleResultTpl))
}

func applyScaleResult(tpl *template.Template, data interface{}) string {
	buf := bytes.NewBuffer(nil)
	if err := tpl.Execute(buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}

type score struct {
	Score int
}

// TODO: 实现一个计算器：基于解析表达式为后缀表达式的函数
func calcScoreWithsRules(rawScores []score, rules []string) []int {
	//return nil
	return []int{10, 23}
}
