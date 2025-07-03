build:
	cargo build --release --target wasm32-unknown-unknown
	wasm2wat ../target/wasm32-unknown-unknown/release/wasm_demo.wasm -o ../demo.wat