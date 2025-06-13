use std::sync::atomic::AtomicUsize;
use std::sync::{mpsc, Arc, Condvar, Mutex};
use std::thread;
use std::time::Duration;

/*
 * 编写一个多线程交替打印的程序，要求如下：
 * 1. 两个线程交替打印，一个线程打印奇数，一个线程打印偶数
 * 2. 打印到100
 * 要求按照顺序输出，不能乱序输出
 */


fn main() {
    // example1(); 这种写法只会先完成一个线程，然后再完成另一个线程：在 rust 中
    // example1();

    // example2(); // 使用 Mutex 和 CondVar 来实现

    // example3(); // 使用 Channel 来实现

    example4(); // 使用 atomic 来实现
}

#[allow(dead_code, unused)]
fn example1() {
    let handle1 = thread::spawn(|| {
        let mut i = 1;
        let start = std::time::Instant::now();
        while i <= 100 {
            println!("thread1: {}", i);
            i += 2;
            // 如果没有 sleep，那么两个线程极大概率会先完成一个线程，然后再完成另一个线程。
            // 这是因为, 打印 100 数大约需要：0.060 ms, 而创建一个线程大概需要 0.240 ms（仅在当前示例中，一般在 微秒 级别）
            // 这里能明显看出在一个线程创建启动之前，可能另一个线程已经执行完成了。所以会出现这种情况。
            // thread::sleep(Duration::from_millis(1));
        }

        // let cost = start.elapsed().as_nanos() as f64 / 1000.0; // 转换为 微秒（0.001 ms）
        // println!("thread1 cost: {} us", cost);
    });

    let handle2 = thread::spawn(|| {
        let mut i = 2;
        while i <= 100 {
            println!("thread2: {}", i);
            i += 2;
            // thread::sleep(Duration::from_millis(1));
        }
    });

    handle1.join().unwrap();
    handle2.join().unwrap();

    // 主线程等待子线程结束
    println!("Main thread finished.");
}

#[allow(dead_code, unused)]
fn example2() {
    // 为了实现两个线程交替执行，可以使用 Mutex 和 CondVar 来实现。
    // Mutex 用于保护共享资源，CondVar 用于线程间的通信。

    let mutex = Arc::new(Mutex::new(1)); // 1 代表打印基数，2 代表打印偶数
    let condvar = Arc::new(Condvar::new());

    let m1 = mutex.clone();
    let c1 = condvar.clone();

    // move 关键字表示 closure 按值捕获，而不是按引用捕获（默认）。
    let handle1 = thread::spawn(move || {
        let mut i = 1;
        while i <= 100 {
            let mut guard = m1.lock().unwrap();
            while *guard != 1 {
                guard = c1.wait(guard).unwrap();
            }

            println!("thread1: {}", i);
            i += 2;

            // switch to print even number
            *guard = 2;
            c1.notify_one();
        }
    });

    let m2 = mutex.clone();
    let c2 = condvar.clone();
    let handle2 = thread::spawn(move || {
        let mut i = 2;
        while i <= 100 {
            let mut guard = m2.lock().unwrap();
            while *guard != 2 {
                // spin and wait
                guard = c2.wait(guard).unwrap();
            }

            println!("thread2: {}", i);
            i += 2;

            // switch to print odd number
            *guard = 1;
            c2.notify_one();
        }
    });

    handle1.join().unwrap();
    handle2.join().unwrap();

    println!("Main thread finished.");
}

#[allow(dead_code, unused)]
fn example3() {
    // 除了使用 Mutex 和 CondVar 来实现，还可以使用 Channel 来实现。
    let (ch1_tx, ch1_rx): (mpsc::Sender<i32>, mpsc::Receiver<i32>) = mpsc::channel(); // channel1 用于打印奇数
    let (ch2_tx, ch2_rx): (mpsc::Sender<i32>, mpsc::Receiver<i32>) = mpsc::channel(); // channel2 用于打印偶数

    let handle1 = thread::spawn(move || {
        let mut i = 0;

        // thread1 从 ch1_rx 中接收数据，然后打印, 再把数据发送到 ch2_tx 中
        // rust 中 ch.recv 是 Result<T, E> 类型，如果 ch 关闭，则返回 RecvError
        for i in ch1_rx {
            if i > 100 {
                break;
            }

            println!("thread1: {}", i);
            // 通知 thread2 可以打印了
            ch2_tx.send(i + 1).unwrap();
        }
    });

    let handle2 = thread::spawn(move || {
        // 启动 thread1
        ch1_tx.send(1).unwrap();

        let mut i = 0;

        // thread2 从 ch2_rx 中接收数据，然后打印, 再把数据发送到 ch1_tx 中
        for i in ch2_rx {
            if i > 100 {
                break;
            }
            println!("thread2: {}", i);
            // 通知 thread1 可以打印了
            ch1_tx.send(i + 1).unwrap();
        }
    });


    handle1.join().unwrap();
    handle2.join().unwrap();

    // 主线程等待子线程结束
    println!("Main thread finished.");
}


#[allow(dead_code, unused)]
fn example4() {
    // 使用 atomic 来实现
    let counter: Arc<AtomicUsize> = Arc::new(AtomicUsize::new(0));

    let counter1 = counter.clone();

    let handle1 = thread::spawn(move || {
        let mut i = 1;

        loop {
            // pub fn compare_exchange(& self, current: usize, new: usize, success: Ordering, failure: Ordering) -> Result<usize, usize>
            // 其中 success 和 failure 分别指示成功和失败时的内存顺序。
            // `内存顺序`: 内存顺序是指在多线程环境中，线程之间如何访问内存。
            // `内存顺序`: 内存顺序有以下几种：
            // 1. Relaxed: 无任何同步约束，可以乱序执行。
            // 2. Release: 保证它之前的操作永远在它之前，它之后的操作可能重排到它之前
            // 3. Acquire: 保证它之后的访问永远在它之后，它之前的操作可能重排到它之后
            // 4. AcqRel: Acquire 和 Release 的组合，同时提供了 Acquire 和 Release 的保证
            // 5. SeqCst: 顺序一致性，像是 AcqRel 的加强版
            match counter1.compare_exchange(
                i - 1,
                i,
                std::sync::atomic::Ordering::SeqCst,
                std::sync::atomic::Ordering::SeqCst,
            ) {
                Ok(_) => {
                    println!("thread1: {}", i);
                    i += 2;
                }
                Err(e) => {
                    // println!("thread1: CAS failed, current: {}, not expected: {}", e, i - 1);
                    thread::sleep(Duration::from_millis(1));
                    // 让出 CPU 时间片，让其他线程有机会执行
                    // thread::yield_now();
                }
            }

            if i > 100 {
                break;
            }
        }

        println!("thread1 finished.");
    });

    let counter2 = counter.clone();
    let handle2 = thread::spawn(move || {
        let mut i = 2;
        loop {
            match counter2.compare_exchange(
                i - 1,
                i,
                std::sync::atomic::Ordering::SeqCst,
                std::sync::atomic::Ordering::SeqCst,
            ) {
                Ok(_) => {
                    println!("thread2: {}", i);
                    i += 2;
                }
                Err(e) => {
                    // println!("thread2: CAS failed, current: {}, not expected: {}", e, i - 1);
                    thread::sleep(Duration::from_millis(1));
                    // 让出 CPU 时间片，让其他线程有机会执行
                    // thread::yield_now();
                }
            }

            if i > 100 {
                break;
            }
        }

        println!("thread2 finished.");
    });

    handle1.join().unwrap();
    handle2.join().unwrap();

    // 主线程等待子线程结束
    println!("Main thread finished.");
}