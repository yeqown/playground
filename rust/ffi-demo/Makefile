.PHONY: all build-rust build-c run-rust run-c clean

all: build-rust build-c

# Build the Rust library
build-rust:
	cargo build

# Build the C program that uses the Rust library
build-c: build-rust
	cc -o c_program c_src/main.c -L../target/debug -lrust_lib -Wl,-rpath,../target/debug

run-c: build-c
	./c_program

clean:
	cargo clean
	rm -f c_program
