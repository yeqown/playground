package calculator_with_tpl

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/pkg/errors"
)

var (
	scaleResultTpl = `
{{ $d1 := index .scores 0 }}
{{ $d2 := index .scores 1 }}
{{ define "result1" }} 多动症 {{ end }}
{{ define "result2" }} 密集恐惧症 {{ end }}
{{ define "result3" }} 幽闭恐惧症 {{ end }}
{{if ge $d1 20.0}}得分：{{$d1}},{{$d2}}, {{- template "result1" }}
{{ else if ge $d2 20.0 }}得分：{{$d1}},{{$d2}}, {{- template "result2" }}
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

type IScorer interface {
	Score() float64
}

type score struct {
	score float64
}

func (s score) Score() float64 {
	return s.score
}

// TODO: 实现一个计算器：基于解析表达式为后缀表达式的函数
func calcScoreWithsRules(rawScores []IScorer, rules []string) ([]float64, error) {
	scores := make([]float64, len(rules))

	for idx, ruleExpr := range rules {
		fn := parseSpecialExprAsPolishWrapper(rawScores)
		score, err := Calculate(ruleExpr, fn)
		if err != nil {
			return nil, errors.Wrap(err, "calcScoreWithsRules.Calculate")
		}

		scores[idx] = score
	}

	return scores, nil
}

// 从中缀表达式转后缀表达式子
// 支持浮点数 和 多位数 和 特殊表达式解析和赋值
// 如：
// @rawScores = [1,2,3,4]
// @expr = q1 * 2 + q2 * 3 + (q3 + q4) / 2
// 转换为后缀表达式时，q1 => 1, q2 => 2 等等
//
func parseSpecialExprAsPolishWrapper(rawScores []IScorer) ParseExprAsPolishFunc {
	return func(expr string) (string, error) {
		var output = make([]string, 0, 32)
		var opStack = newStack()

		var (
			// 用于读取数值和特殊表达式
			preNum           string
			flagNumContinue  bool
			flagSpecialStart bool
		)

		outputNum := func() error {
			// 检查数字是否为空
			if preNum != "" {
				// 改数值是特殊字符开头的
				if flagSpecialStart {
					idx, err := strconv.Atoi(preNum)
					if err != nil {
						return errors.Wrap(err, "parseSpecialExprAsPolishWrapper.parse special idx")
					}
					score := rawScores[idx-1].Score()
					preNum = strconv.FormatFloat(score, 'f', -1, 64)
				}

				output = append(output, preNum)

				preNum = ""
				flagNumContinue = false
				flagSpecialStart = false
			}
			return nil
		}

		for _, r := range expr {
			// fmt.Println("r=", string(r), "output=", output)
			if r == ' ' {
				if err := outputNum(); err != nil {
					return "", err
				}
				continue
			}

			// 这里特殊处理, 处理数字
			if unicode.IsDigit(r) {
				// True: 数字
				flagNumContinue = true
				preNum += string(r)
				continue
			} else if r == '.' && flagNumContinue {
				// True: 小数点
				preNum += string(r)
				continue
			} else if r == 'q' && !flagSpecialStart {
				flagSpecialStart = true
				continue
			}

			if err := outputNum(); err != nil {
				return "", err
			}

			// 操作符
			if isOp(r) {
				if r == ')' {
					// 遇到 ')'则开始弹出 opStack，直到遇见'('
					for !opStack.Empty() {
						peak := opStack.Pop().(rune)
						if peak == '(' {
							break
						}
						output = append(output, string(peak))
					}

					// fmt.Printf("遇见 ')' 处理后 stack=%s\n", runeHelper(opStack))
					continue
				}

				// 除了 ')' 以外的操作符，opStack 一直弹出直到 栈顶元素优先级 小于 当前操作符
				for !opStack.Empty() && !isOpLower(opStack.Peak().(rune), r) {
					// 栈不为空 且 栈顶元素优先级更高（或相等）
					peak := opStack.Peak().(rune)
					if peak == '(' {
						break
					}
					opStack.Pop()
					output = append(output, string(peak))
				}
				opStack.Push(r)
				// fmt.Printf("遇见 '%s' 处理后，stack=%s\n", string(r), runeHelper(opStack))
			}
			// TODO: 不能处理的字符，忽略或者报错
		}

		// FIX: 以数字结尾的表达式
		if flagNumContinue {
			output = append(output, preNum)
		}

		// 表达式遍历完成，弹出所有操作符
		for !opStack.Empty() {
			r := opStack.Pop().(rune)
			output = append(output, string(r))
		}

		return strings.Join(output, ","), nil
	}
}
