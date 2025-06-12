# Rust 练习环境

这是一个简化的 Rust 练习环境，用于创建和运行多个独立的 Rust 文件，并共享库代码。

## 目录结构

```
rust/
├── Cargo.toml         # 项目配置文件
├── .gitignore         # Git 忽略文件
└── src/               # 源代码目录
    ├── lib.rs         # 库代码的主入口点
    ├── math.rs        # 数学相关的库模块
    ├── main.rs        # 主项目入口点
    └── bin/           # 独立可执行文件目录
        ├── hello_world.rs
        ├── fibonacci.rs
        ├── simple_calculator.rs
        ├── use_lib_example.rs
        └── math_example.rs
```

## 如何使用

### 运行独立的 Rust 文件

每个放在 `src/bin` 目录下的 Rust 文件都可以作为独立程序运行：

```bash
# 运行 hello_world.rs
cd rust
cargo run --bin hello_world

# 运行使用库代码的示例
cargo run --bin use_lib_example
cargo run --bin math_example
```

### 创建新的 Rust 文件

只需在 `src/bin` 目录下创建新的 `.rs` 文件，确保文件中包含 `main()` 函数：

```bash
# 创建新文件
touch src/bin/my_new_script.rs

# 编辑文件，添加 main() 函数
# ...

# 运行
cargo run --bin my_new_script
```

### 使用库代码

在任何 bin 目录下的文件中，你可以通过 crate 名称（在 Cargo.toml 中定义的 name）来引用库代码：

```rust
// 使用库中的函数
let result = rust_playground::add(10, 20);

// 使用库中的模块
let is_prime = rust_playground::math::is_prime(17);
```

### 添加新的库模块

1. 在 `src` 目录下创建新的 `.rs` 文件
2. 在 `lib.rs` 中使用 `pub mod` 声明该模块

例如，创建 `src/strings.rs` 后，在 `lib.rs` 中添加：

```rust
pub mod strings;
```

### 运行测试

```bash
# 运行所有测试
cargo test

# 运行特定模块的测试
cargo test --lib math
```

### 添加依赖

如果需要添加依赖，只需在 `Cargo.toml` 文件中的 `[dependencies]` 部分添加即可，所有的 bin 文件和库代码都可以使用这些依赖。

例如：

```toml
[dependencies]
rand = "0.8.5"
serde = { version = "1.0", features = ["derive"] }
```
