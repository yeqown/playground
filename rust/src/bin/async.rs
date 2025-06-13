/*
 * 实践 Rust 中的 Future 编程的 和 async/await 编程
 */

use futures::executor::block_on;

#[async_std::main]
async fn main() {
    example1(); // 异步任务的执行

    example2(); // 多个异步任务的组合
}

fn example1() {
    let fut = doing("coding");
    // block_on 的作用是：让当前线程阻塞，直到 future 完成。
    block_on(fut);
}

async fn doing(work: &str) {
    println!("async doing: {}", work);
}

async fn ordered_doing(work1: &str, work2: &str) {
    doing(work1).await;
    doing(work2).await;
}

async fn multi_future() {
    let fut1 = ordered_doing("go home", "cooking");
    let fut2 = doing("listening to music");

    // 可以并发的处理和等待多个 future
    futures::join!(fut1, fut2);
}

fn example2() {
    block_on(multi_future());
}