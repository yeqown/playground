/*
 * 宏编程是通过一种代码来生成代码的编程方式。
 * 于 元编程 类似，都是为了减少重复代码的编写。宏编程有以下特点：
 * 1. 可变参数
 * 2. 编译前展开
 * 3. 宏自身难以理解和维护
 *
 * 宏分为两种：
 * - `声明式宏`：宏也是将一个值跟对应的 模式 进行匹配，且该模式会与特定的代码相关联。如：
 *   ```
 *   let a = 1;
 *   println!("a = {}", a);
 *   ```
 * - `过程式宏`：使用 proc_macro 来定义宏（derive，属性宏，函数宏）。过程宏需要先编译才能使用。它不是用来替换代码，而是生成代码,如：
 *   derive 宏：
 *   ```
 *   #[derive(Debug)]
 *   struct Person {
 *       name: String,
 *       age: u32,
 *   }
 *   ```

 *  属性宏：
 *  ```
 *  #[route("/hello")]
 *  fn hello() {
 *     println!("Hello, world!");
 *  }
 *
 *  类函数宏：
 *  ```
 *  #[proc_macro]
 *  fn sql(input: TokenStream) -> TokenStream {}
 *
 *  let sql = sql!("SELECT * FROM users");
 *  ```
 */

use rust_playground_macros::*;

#[derive(Debug)]
struct User {
    name: String,
    age: u32,
}

// 使用路由宏
#[log]
fn get_users(user_id: u64) -> User {
    println!("fetching by user_id: {}", user_id);

    User {
        name: format!("Robot-{}", user_id),
        age: 20,
    }
}

fn main() {
    // 声明式宏
    println!("Hello, macro!");

    // 过程式宏 - derive 宏
    let u = User {
        name: String::from("Alice"),
        age: 20,
    };
    println!("macro derive: {:?}, name={:?}, age={:?}", u, u.name, u.age);

    // 过程式宏 - 属性宏
    get_users(10u64);

    // 过程式宏 - 类函数宏
    let query = sql!("SELECT * FROM users WHERE active = true");
    println!("Query string: {}", query);
}


