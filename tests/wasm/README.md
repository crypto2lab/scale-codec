### WASM

Check the examples at:

- `testdata/main.go`: here we have a golang code that will be compiled to wasm, it has a simple `addTwoNumbers` method that is exported, then this method receives two uint64 arguments, the argument is to locate where in the memory is written the scale encoded bytes. After reading the bytes from the linear memory we sum them up and scale encode the result and store it in the linear memory again returning the pointer + size to the caller.

- `src/main.rs`: here is where we instantiate the wasm blob, encode and write in the memory the arguments, call the exported `addTwoNumbers` function and check the output

Compiling the golang to wasm:

- Make sure to have (tinygo)[#https://tinygo.org/getting-started/] installed

```sh
cd tests/wasm/testdata; tinygo build -o main.wasm -scheduler=none -no-debug -target=wasi main.go
```

- To run the rust code:

```sh
cd tests/wasm; cargo run
```

The output should be:

```sh
Compiling module...
Starting `tokio` runtime...
Creating `WasiEnv`...
Instantiating module...
expected: 3, from wasm: 3
```
