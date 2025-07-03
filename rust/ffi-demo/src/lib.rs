use std::os::raw::c_int;

// FFI bindings to C functions
#[link(name = "c_functions")]
extern "C" {
    fn add_numbers(a: c_int, b: c_int) -> c_int;
    fn call_rust_function(
        rust_fn: extern "C" fn(c_int, c_int) -> c_int,
        a: c_int,
        b: c_int,
    ) -> c_int;
}

// Rust function that will be called from C
#[no_mangle]
pub extern "C" fn multiply_numbers(a: c_int, b: c_int) -> c_int {
    println!("Rust function 'multiply_numbers' called with {} and {}", a, b);
    a * b
}

// Safe wrapper for the C add_numbers function
pub fn safe_add_numbers(a: i32, b: i32) -> i32 {
    unsafe { add_numbers(a, b) }
}

// Safe wrapper to call C function that calls back to Rust
pub fn safe_call_rust_function(a: i32, b: i32) -> i32 {
    unsafe { call_rust_function(multiply_numbers, a, b) }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_c_add() {
        assert_eq!(safe_add_numbers(5, 7), 12);
    }

    #[test]
    fn test_rust_multiply_via_c() {
        assert_eq!(safe_call_rust_function(5, 7), 35);
    }
}
