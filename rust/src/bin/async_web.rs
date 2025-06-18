use async_std::net::TcpListener;
use async_std::net::TcpStream;
use async_std::prelude::*;
use async_std::task;
use async_std::task::spawn;
use futures::stream::StreamExt;
use std::fs;
// use std::io::prelude::*;
// use std::net::TcpListener;
// use std::net::TcpStream;
use std::path::Path;
use std::time::Duration;

/// 注意 ⚠️ 使用浏览器测试 sleep 接口时会出现阻塞的现象：
/// 有些浏览器会对同一 URL 的请求进行 pending 而不是立即发送。

#[async_std::main]
async fn main() {
    // 监听本地端口 7878 ，等待 TCP 连接的建立
    let listener = TcpListener::bind("127.0.0.1:7878").await.unwrap();

    listener
        .incoming()
        .for_each_concurrent(/* Futures limit */Some(2), |stream| async move {
            let stream = stream.unwrap();

            #[allow(unused)]
            /// async_std::task::spawn 可以将一个 Future 放入到任务队列中，
            /// 并立即返回，这样可以继续处理下一个连接。
            /// 不会阻塞当前线程，而是在后台执行 Future。
            ///
            /// async_std 支持多线程执行。类似于 golang 的 GPM 模型。
            spawn(handle_connection(stream));

            // 这种写法可以实现并发，但不是并行，同时受到 Futures limit 的限制
            // handle_connection(stream).await;
        })
        .await;
}

async fn handle_connection(mut stream: TcpStream) {
    // 从连接中顺序读取 1024 字节数据
    let mut buffer = [0; 1024];
    let _ = stream.read(&mut buffer).await.unwrap();

    let get = b"GET / HTTP/1.1\r\n";
    let sleep = b"GET /sleep HTTP/1.1\r\n";


    // 处理HTTP协议头，若不符合则返回404和对应的 `html` 文件
    let (status_line, filename) = if buffer.starts_with(get) {
        ("HTTP/1.1 200 OK\r\n\r\n", "hello.html")
    } else if buffer.starts_with(sleep) {
        task::sleep(Duration::from_secs(5)).await;
        ("HTTP/1.1 200 OK\r\n\r\n", "hello.html")
    } else {
        ("HTTP/1.1 404 NOT FOUND\r\n\r\n", "404.html")
    };

    let filepath = Path::new("assets").join(filename);
    let contents = fs::read_to_string(filepath).unwrap();

    // 将回复内容写入连接缓存中
    let response = format!("{status_line}{contents}");
    stream.write(response.as_bytes()).await.unwrap();
    // 使用 flush 将缓存中的内容发送到客户端
    stream.flush().await.unwrap();
}