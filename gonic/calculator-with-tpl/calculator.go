// https://blog.csdn.net/ssjhust123/article/details/8001651

package calculator_with_tpl

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// opPriority ...
var opPriority = map[rune]int{
	'+': 1, '-': 1,
	'*': 2, '/': 2,
	'(': 3, ')': 3,
}

// 从中缀表达式转后缀表达式子
type ParseExprAsPolishFunc func(expr string) (string, error)

func Calculate(expr string, fn ParseExprAsPolishFunc) (float64, error) {
	polishExpr, err := fn(expr)
	if err != nil {
		return 0.0, errors.Wrap(err, "Calculate.ParseExprAsPolishFunc")
	}
	return calcWithPolish(polishExpr), nil
}

// calcWithPolish 根据后缀表达式求值
func calcWithPolish(polishExpr string) float64 {
	outs := strings.Split(polishExpr, ",")
	s := newStack()

	for _, v := range outs {
		num, err := strconv.ParseFloat(v, 64)
		if err == nil {
			// true: 如果可以转成数字
			s.Push(num)
			continue
		}

		// fmt.Printf("计算前: stack=%s\n", float64Helper(s))
		op1 := s.Pop().(float64)
		op2 := s.Pop().(float64)
		switch v {
		case "+":
			s.Push(op2 + op1)
		case "-":
			s.Push(op2 - op1)
		case "*":
			s.Push(op2 * op1)
		case "/":
			s.Push(op2 / op1)
		}
		// fmt.Printf("计算 %.2f %s %.2f: stack=%s\n", op1, v, op2, float64Helper(s))
	}
	return s.Pop().(float64)
}

// 从中缀表达式转后缀表达式子
// 支持浮点数和多位数
func defaultParseExprAsPolish(expr string) (string, error) {
	var output = make([]string, 0, 32)
	var opStack = newStack()

	var (
		// 用于读取数值
		preNum          string
		flagNumContinue bool
	)

	for _, r := range expr {
		// fmt.Println("r=", string(r), "output=", output)
		if r == ' ' {
			if preNum != "" {
				output = append(output, preNum)
			}

			preNum = ""
			flagNumContinue = false
			continue
		}

		// 处理数字
		if unicode.IsDigit(r) {
			// True: 数字
			flagNumContinue = true
			preNum += string(r)
			continue
		} else if r == '.' && flagNumContinue {
			// True: 小数点
			preNum += string(r)
			continue
		}

		// 数字表达结束，将其输出，并重置标志位
		if flagNumContinue {
			output = append(output, preNum)
			preNum = ""
			flagNumContinue = false
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

// if c (rune) is a char in ops return true
func isOp(c rune) bool {
	_, ok := opPriority[c]
	return ok
}

// isOpLower means to do
// return r1 < r2
func isOpLower(r1, r2 rune) bool {
	return opPriority[r1] < opPriority[r2]
}
