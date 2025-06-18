use std::sync::Arc;
use std::thread;

/*
 * Send 和 Sync Trait 是 Rust 中并发安全的两个重要 标记 特性。
 * 实现 Send 的类型可以在线程间安全的传递所有权。
 * 实现 Sync 的类型可以在线程间安全的共享（引用）。
 *
 * Rust 中大部分类型都实现了 Send 和 Sync, 除了：
 * 1. 原始指针（裸指针）
 * 2. UnsafeCell, Cell 和 RefCell, 没有实现 Sync。（这几个类型是用来）
 * 3. Rc 都没有实现
 * 4. 自定义类型中包含了未实现 Send 和 Sync 的类型
 */

fn main() {
    // 使用 Arc 来实现线程安全的共享
    let v = Arc::new(5);
    let t = thread::spawn(move || {
        println!("{}", v);
    });
    t.join().unwrap();

    // 自定义实现 Send
    let p = MyBox(5 as *mut u8);
    let t = thread::spawn(move || {
        println!("{:?}", p);
    });

    t.join().unwrap();
}

#[derive(Debug)]
struct MyBox(*mut u8);

unsafe impl Send for MyBox {}

#[cfg(ignore)]
#[allow(dead_code)]
fn version1() {
    use std::rc::Rc;

    // the trait `Send` is not implemented for `Rc<i32>`
    // Rc 未实现 Send trait，因此不能在线程间传递
    let v = Rc::new(5);
    let t = thread::spawn(move || {
        println!("{}", v);
    });

    t.join().unwrap();
}
