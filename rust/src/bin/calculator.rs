/*
 * This is a calculator to practice compound condition control,
 * built-in types, match pattern and high-level data structures.
 *
 * Expect input: ( 1 + 2 ) + 3 / 3
 */
use rust_playground::stack::Stack;
use std::fmt::Debug;
use thiserror::Error;

// Custom Error for calculator
#[derive(Debug, PartialEq, Error)]
enum MyError {
    #[error("Invalid Expression: {0}")]
    InvalidExpression(String),
    #[error("Division by Zero: {0}")]
    DivisionByZero(String),
}

trait Calculator {
    fn eval(expr: &str) -> Result<i64, MyError>;
}

#[derive(Debug, PartialEq)]
enum ExprToken {
    OPERATOR(char),
    OPERAND(i64),
}

#[derive(Default)]
struct CalculatorBasedStack {}

impl CalculatorBasedStack {
    /// Parse the input expression into a vector of tokens.
    ///
    /// # Arguments
    ///
    /// * `expr` - The input expression to parse.
    ///
    /// # Returns
    ///
    /// A vector of tokens.
    ///
    /// # Examples
    ///
    /// ```
    /// let expr = "1 + 2 + 3 / 3";
    /// let tokens = parse(expr);
    /// ```
    fn parse(&self, expr: &str) -> Vec<ExprToken> {
        let mut tokens = Vec::new();
        let mut chars = expr.chars().peekable();
        let mut num_buffer = String::new();

        while let Some(c) = chars.next() {
            if c.is_whitespace() {
                continue;
            }

            if c.is_digit(10) {
                num_buffer.push(c);
                while let Some(&next) = chars.peek() {
                    if next.is_digit(10) {
                        num_buffer.push(next);
                        chars.next();
                    } else {
                        break;
                    }
                }
                if let Ok(num) = num_buffer.parse::<i64>() {
                    tokens.push(ExprToken::OPERAND(num));
                }
                num_buffer.clear();
            } else {
                tokens.push(ExprToken::OPERATOR(c));
            }
        }

        tokens
    }

    fn apply_operator(operators: &mut Stack<char>, operands: &mut Stack<i64>) -> Result<(), MyError> {
        if let (Some(op), Some(rhs), Some(lhs)) = (operators.pop(), operands.pop(), operands.pop()) {
            let result = match op {
                '+' => lhs + rhs,
                '-' => lhs - rhs,
                '*' => lhs * rhs,
                '/' => {
                    if rhs == 0 {
                        return Err(MyError::DivisionByZero("Division by zero".to_string()));
                    }
                    lhs / rhs
                }
                _ => return Err(MyError::InvalidExpression(format!("Unknown operator: {}", op))),
            };
            operands.push(result);
            Ok(())
        } else {
            Err(MyError::InvalidExpression("Invalid expression".to_string()))
        }
    }

    fn precedence(op: char) -> u8 {
        match op {
            '+' | '-' => 1,
            '*' | '/' => 2,
            _ => 0,
        }
    }
}

impl Calculator for CalculatorBasedStack {
    fn eval(expr: &str) -> Result<i64, MyError> {
        let mut operators = Stack::new();
        let mut operands = Stack::new();
        let tokens = CalculatorBasedStack::default().parse(expr);

        #[allow(unused_doc_comments)]
        /// 遍历表达式中的每个token
        /// 如果是操作数，直接压入操作数栈
        /// 如果是运算符，则根据优先级进行处理

        for token in tokens {
            match token {
                ExprToken::OPERAND(num) => operands.push(num),
                ExprToken::OPERATOR(op) => {
                    match op {
                        '(' => operators.push(op), // 左括号直接入栈
                        ')' => { // 右括号，则需要弹出栈顶的运算符并运算，直到遇到左括号
                            while let Some(&top_op) = operators.peek() {
                                if top_op != '(' {
                                    CalculatorBasedStack::apply_operator(&mut operators, &mut operands)?;
                                } else {
                                    break;
                                }
                            }
                            // 弹出左括号
                            if let Some('(') = operators.pop() {
                                // 正常处理
                            } else {
                                return Err(MyError::InvalidExpression("Mismatched parentheses".to_string()));
                            }
                        }
                        _ => {
                            // 非括号的运算符，需要比较优先级，比较当前操作符和栈顶操作符的优先级
                            // 如果 栈顶操作符 的优先级则
                            while let Some(&top_op) = operators.peek() {
                                if CalculatorBasedStack::precedence(top_op) >= CalculatorBasedStack::precedence(op) && top_op != '(' {
                                    CalculatorBasedStack::apply_operator(&mut operators, &mut operands)?;
                                } else {
                                    break;
                                }
                            }
                            operators.push(op);
                        }
                    }
                }
            }
        }

        while operators.peek().is_some() {
            if let Some(&top_op) = operators.peek() {
                if top_op == '(' {
                    return Err(MyError::InvalidExpression("Mismatched parentheses".to_string()));
                }
            }
            CalculatorBasedStack::apply_operator(&mut operators, &mut operands)?;
        }

        if let Some(result) = operands.pop() {
            if operands.is_empty() {
                Ok(result)
            } else {
                Err(MyError::InvalidExpression("Invalid expression".to_string()))
            }
        } else {
            Err(MyError::InvalidExpression("Invalid expression".to_string()))
        }
    }
}

#[test]
fn test_calculator_parse() {
    let expr = "1 + 2 + 3 / 3";
    let tokens = CalculatorBasedStack::default().parse(expr);
    assert_eq!(
        tokens,
        vec![
            ExprToken::OPERAND(1),
            ExprToken::OPERATOR('+'),
            ExprToken::OPERAND(2),
            ExprToken::OPERATOR('+'),
            ExprToken::OPERAND(3),
            ExprToken::OPERATOR('/'),
            ExprToken::OPERAND(3)
        ]
    )
}

#[test]
fn test_calculator_eval() {
    let expr = "1 + 24 + 3 / 3";
    let result = CalculatorBasedStack::eval(expr);
    assert_eq!(result, Ok(26));
}

#[test]
fn test_calculator_eval_with_parentheses() {
    let expr = "(1 + 2) * (3 + 3)";
    let result = CalculatorBasedStack::eval(expr);
    assert_eq!(result, Ok(18));
}

fn main() {
    let expr = "1 + 2 + 3 / 3";
    match CalculatorBasedStack::eval(expr) {
        Ok(result) => println!("计算结果: {}", result),
        Err(e) => println!("计算出错: {}", e),
    }
}