// https://blog.csdn.net/ssjhust123/article/details/8001651

package tpl

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// opPriority ...
var opPriority = map[rune]int{
	43: 1, // +
	45: 1, // -
	42: 2, // *
	47: 2, // /
	40: 3, // (
	41: 3, // )
}

func Calculate(expr string) float32 {
	polishExpr := parseExprAsPolishV2(expr)
	return calcWithPolish(polishExpr)
}

// calcWithPolish ...
// TODO: 修改这部分
func calcWithPolish(polishExpr string) float32 {
	outs := strings.Split(polishExpr, ",")
	s := newStack()

	for i := len(outs) - 1; i >= 0; i-- {
		num, _ := strconv.ParseFloat(outs[i], 64)
		if num < 0 {
			// 操作符
			num = 0 - num
			switch rune(num) {
			case '+':
				s.Push(s.Pop().(int) + s.Pop().(int))
			case '-':
				s.Push(s.Pop().(int) - s.Pop().(int))
			case '*':
				s.Push(s.Pop().(int) * s.Pop().(int))
			case '/':
				s.Push(int(s.Pop().(int) / s.Pop().(int)))
			default:
				fmt.Printf("Error: invalid op = %v", num)
			}
		} else {
			s.Push(int(num))
		}
	}

	return s.Pop().(float32)
}

func parseExprAsPolishV2(expr string) string {
	var output = make([]string, 0, 32)
	var opStack = newStack()

	var (
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

		if isOp(r) {
			if r == ')' {
				// 遇到 ')' 弹出栈，直到遇见'('
				for !opStack.Empty() {
					peak := opStack.Pop().(rune)
					output = append(output, string(peak))
					if peak == '(' {
						break
					}
				}
				continue
			}

			// output = append(output, string(r))
			for !opStack.Empty() && isOpLower(r, opStack.Peak().(rune)) {
				// 栈不为空 且 栈顶元素优先级更高（或相等）
				peak := opStack.Pop().(rune)
				if peak == '(' {
					break
				}
				output = append(output, string(peak))
			}
			opStack.Push(r)
		}
		// TODO: 不能处理的字符，忽略或者报错
	}

	if flagNumContinue {
		output = append(output, preNum)
	}

	// 表达式遍历完成，弹出所有操作符
	for !opStack.Empty() {
		r := opStack.Pop().(rune)
		output = append(output, string(r))
	}

	return strings.Join(output, ",")
}

var (
	_ops = map[rune]bool{
		'+': true,
		'-': true,
		'*': true,
		'/': true,
		'(': true,
		')': true,
	}
)

//func isNumber(c rune) bool {
//	_, ok := _numbers[c]
//	return ok
//}

func isOp(c rune) bool {
	_, ok := _ops[c]
	return ok
}

//
//func isParenthesis(c rune) bool {
//	return c == '(' || c == ')'
//}
//
//func isRightParenthesis(c rune) bool {
//	return c == ')'
//}

// isOpLower means to do
// return r1 < r2
func isOpLower(r1, r2 rune) bool {
	return opPriority[r1] <= opPriority[r2]
}
