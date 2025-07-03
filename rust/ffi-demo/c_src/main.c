#include <stdio.h>
#include "c_functions.h"

// Import the Rust function
extern int multiply_numbers(int a, int b);

int main() {
    printf("C Program Using Rust Library\n");
    printf("============================\n");

    // Call Rust function directly
    int a = 5;
    int b = 8;
    int result = multiply_numbers(a, b);
    printf("Result from Rust multiply_numbers(%d, %d): %d\n", a, b, result);

    return 0;
}