use std::io::Read;
use std::{borrow::BorrowMut, error::Error};

use parity_scale_codec::{self, Decode, Encode};
use tokio;
use wasmer::{imports, Bytes, Instance, Module, Pages, Store, TypedFunction, WasmPtr};
use wasmer_wasix::virtual_net::BytesMut;
use wasmer_wasix::{Pipe, WasiEnv};

fn main() -> Result<(), Box<dyn Error>> {
    let wasm_path = "testdata/main.wasm";
    // Let's declare the Wasm module with the text representation.
    let wasm_bytes = std::fs::read(wasm_path)?;

    // Create a Store.
    let mut store = Store::default();

    println!("Compiling module...");
    // Let's compile the Wasm module.
    let module = Module::new(&store, wasm_bytes)?;

    println!("Starting `tokio` runtime...");
    let runtime = tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .unwrap();
    let _guard = runtime.enter();

    println!("Creating `WasiEnv`...");
    // First, we create the `WasiEnv`
    let mut wasi_env = WasiEnv::builder("hello").finalize(&mut store)?;

    println!("Instantiating module...");
    // Let's instantiate the Wasm module.
    let import_object = wasi_env.import_object(&mut store, &module)?;
    let instance = Instance::new(&mut store, &module, &import_object)?;

    // Attach the memory export
    let memory = instance.exports.get_memory("memory")?;
    wasi_env.initialize(&mut store, instance.clone())?;

    // We now have an instance ready to be used.
    //
    // We will start by querying the most intersting information
    // about the memory: its size. There are mainly two ways of getting
    // this:
    // * the size as a number of `Page`s
    // * the size as a number of bytes
    //
    // The size in bytes can be found either by querying its pages or by
    // querying the memory directly.
    println!("Querying memory size...");
    let memory_view = memory.view(&store);
    println!("{}", memory_view.size().0);
    println!("{}", memory_view.size().bytes().0);
    println!("{}", memory_view.data_size());

    let x: u64 = 78;
    let y: u64 = 100;

    let enc_x: Vec<u8> = (x as u64).encode();
    let enc_y = (y as u64).encode();

    let malloc: TypedFunction<i32, WasmPtr<u8>> =
        instance.exports.get_typed_function(&store, "malloc")?;

    let x_ptr = malloc.call(&mut store, enc_x.len() as i32)?;
    let y_ptr = malloc.call(&mut store, enc_y.len() as i32)?;

    let memory_view = memory.view(&store);

    println!("A addr: {}", x_ptr.offset());
    println!("B addr: {}", y_ptr.offset());

    x_ptr
        .slice(&memory_view, enc_x.len() as u32)?
        .write_slice(enc_x.as_slice())?;

    y_ptr
        .slice(&memory_view, enc_y.len() as u32)?
        .write_slice(enc_y.as_slice());

    let add_two_number: TypedFunction<(u64, u64), u64> = instance
        .exports
        .get_typed_function(&store, "addTwoNumbers")?;

    let result = add_two_number.call(
        &mut store,
        pointer_size(x_ptr.offset(), enc_x.len() as u32),
        pointer_size(y_ptr.offset(), enc_y.len() as u32),
    )?;

    println!("{}", result);
    let (resultPtr, resultLen) = split(result);
    println!("ptr: {}, len: {}", resultPtr, resultLen);

    let memory_view = memory.view(&store);

    let result_wasm_ptr: WasmPtr<u8> = WasmPtr::new(resultPtr);

    let binding = result_wasm_ptr
        .slice(&memory_view, resultLen)?
        .read_to_vec()?;
    let mut result_contents: &[u8] = binding.as_slice();

    let from_wasm: u64 = u64::decode(&mut result_contents)?;
    let expected_result = x + y;
    println!("expected: {expected_result}, from wasm: {from_wasm}");
    Ok(())
}

fn pointer_size(ptr: u32, size: u32) -> u64 {
    (ptr as u64) << 32 | (size as u64)
}

fn split(ptr_size: u64) -> (u32, u32) {
    ((ptr_size >> 32) as u32, (ptr_size & 0xFFFFFFFF) as u32)
}
