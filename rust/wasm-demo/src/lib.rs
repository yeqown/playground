#![no_std] // 禁用标准库

// 定义一个外部函数，用于打印消息
unsafe extern "C" {
    fn print(result: u64);
}

#[unsafe(no_mangle)]
pub extern "C" fn add(left: u64, right: u64) -> u64 {
    let r = left + right;
    unsafe {
        print(r);
    }
    r
}

#[unsafe(no_mangle)]
pub extern "C" fn subtract(left: u32, right: u32) -> u32 {
    let r = left - right;
    unsafe {
        print(r as u64);
    }
    r
}

#[cfg(not(test))] // Ensures this handler is not used during tests if you have them
#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    // In a no_std environment, the simplest panic handler is an infinite loop.
    // This effectively "traps" the Wasm execution, preventing further execution
    // which is the desired behavior for unrecoverable errors in Wasm.
    loop {}
}