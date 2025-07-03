#ifndef C_FUNCTIONS_H
#define C_FUNCTIONS_H

#ifdef __cplusplus
extern "C" {
#endif

// C function that will be called from Rust
int add_numbers(int a, int b);

// Function pointer type for the Rust function we'll call from C
typedef int (*RustFunctionType)(int, int);

// C function that calls a Rust function (passed as function pointer)
int call_rust_function(RustFunctionType rust_fn, int a, int b);

#ifdef __cplusplus
}
#endif

#endif // C_FUNCTIONS_H
