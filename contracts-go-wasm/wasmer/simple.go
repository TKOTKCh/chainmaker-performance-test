package main

import (
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
	"io/ioutil"
)

func main() {
	methodName := "runtime_type"
	wasmBytes, err := ioutil.ReadFile("../compute/compute-tinygo.wasm")
	if err != nil {
		panic(fmt.Sprintf("Error reading WASM file: %v", err))
	}

	store := wasmer.NewStore(wasmer.NewEngine())
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		panic(fmt.Sprintf("Error creating WASM module: %v", err))
	}

	wasiEnv, err := wasmer.NewWasiStateBuilder("wasi-program").
		Finalize()
	if err != nil {
		panic(fmt.Sprintf("Error creating WASI environment: %v", err))
	}

	importObject, err := wasiEnv.GenerateImportObject(store, module)
	if err != nil {
		panic(fmt.Sprintf("Error generating import object: %v", err))
	}

	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		panic(fmt.Sprintf("Error instantiating module: %v", err))
	}

	// 获取并执行 WASI start 函数
	start, err := instance.Exports.GetWasiStartFunction()
	if err != nil {
		panic(fmt.Sprintf("Error getting WASI start function: %v", err))
	}
	start()
	// 获取并执行 method 函数
	method, err := instance.Exports.GetFunction(methodName)
	if err != nil {
		panic(fmt.Sprintf("Error getting %s function: %v", methodName, err))
	}

	result, err := method()
	if err != nil {
		panic(fmt.Sprintf("Error calling %s function: %v", methodName, err))
	}

	fmt.Println(result) // 输出 3
	//wasmBytes, _ := os.ReadFile("simple.wasm")
	//
	//engine := wasmer.NewEngine()
	//store := wasmer.NewStore(engine)
	//
	//// Compiles the module
	//module, err := wasmer.NewModule(store, wasmBytes)
	//check(err)
	//// Instantiates the module
	//importObject := wasmer.NewImportObject()
	//instance, err := wasmer.NewInstance(module, importObject)
	//check(err)
	//// Gets the `method` exported function from the WebAssembly instance.
	//method, _ := instance.Exports.GetFunction(methodName)
	//
	//// Calls that exported function with Go standard values. The WebAssembly
	//// types are inferred and values are casted automatically.
	//result, _ := method(5, 37)
	//
	//fmt.Println(result) // 42!
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
