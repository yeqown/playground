fn main() {
    // Compile C code
    cc::Build::new()
        .file("c_src/c_functions.c")
        .compile("c_functions");
    
    // Tell cargo to invalidate the built crate whenever the C source changes
    println!("cargo:rerun-if-changed=c_src/c_functions.c");
    println!("cargo:rerun-if-changed=c_src/c_functions.h");
}
