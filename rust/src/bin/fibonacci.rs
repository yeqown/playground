fn fibonacci(n: u32) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => fibonacci(n - 1) + fibonacci(n - 2),
    }
}

fn main() {
    println!("Fibonacci sequence:");
    for i in 0..10 {
        println!("fibonacci({}) = {}", i, fibonacci(i));
    }
}
