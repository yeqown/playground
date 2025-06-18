/*
 * Closure 闭包也就是 python 中的 lambda 函数, go 中的匿名函数
 * 它可以被赋值给变量，也可以作为参数传递给函数。大体上和函数类似，但是它可以捕获调用者作用域中值
 * 1. Rust 中闭包的语法：
 * |param1, param2,...| -> 返回类型 {
        语句1;
        语句2;
        返回表达式
    }
 *
 * 2. 闭包也支持类型推导，但是当编译器推导出一个类型后，就不能再改变它的类型。如下：
 * let echo = |x| x;
 * echo(1);
 * echo("hello"); // 错误：类型不匹配
 *
 * 3. 闭包也支持类型约束，参见 closure_limited_syntax 函数签名
 * 4. 闭包捕获变量捕获变量有三种途径, 分别对应了三种约束
 *   1. 转移所有权 -> FnOnce Trait
 *   2. 可变借用 -> FnMut Trait
 *   3. 不可变借用 -> Fn Trait
 */


fn main() {
    let x = 5u32;
    let sum = |y: u32| -> u32 {
        x + y
    };

    assert_eq!(sum(3), 8);


    // 闭包约束
    let add_10 = |z: u32| -> u32 {
        z + 10
    };

    let mut var = 10;
    let double = |z: u32| -> u32 {
        var = z;
        z + var
    };

    assert_eq!(closure_trait_fn_once(10, add_10), 20);

    assert_eq!(closure_trait_fn_add_twice(10, add_10), 40);

    assert_eq!(closure_trait_mut_fn(5, double), 10);
    assert_eq!(var, 5);
    println!("var: {}", var);
}

fn closure_trait_mut_fn<F>(source: u32, mut add_func: F) -> u32
where
    F: FnMut(u32) -> u32, // 闭包的类型约束
{
    add_func(source)
}

fn closure_trait_fn_once<F>(source: u32, add_func: F) -> u32
where
    F: FnOnce(u32) -> u32, // 闭包的类型约束
{
    let result = add_func(source);
    // result + add_func(source) // `FnOnce` closures can only be called once

    #[allow(clippy::let_and_return)]
    result
}

fn closure_trait_fn_add_twice<F>(source: u32, add_func: F) -> u32
where
    F: Fn(u32) -> u32, // 闭包的类型约束, 不可变借用
{
    let result = add_func(source);
    result + add_func(source)
}