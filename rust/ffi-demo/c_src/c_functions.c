#include "c_functions.h"
#include <stdio.h>

// Implementation of the C function that will be called from Rust
int add_numbers(int a, int b) {
    printf("C function 'add_numbers' called with %d and %d\n", a, b);
    return a + b;
}

// Implementation of the C function that calls a Rust function
int call_rust_function(RustFunctionType rust_fn, int a, int b) {
    printf("C function 'call_rust_function' calling Rust function with %d and %d\n", a, b);
    return rust_fn(a, b);
}
