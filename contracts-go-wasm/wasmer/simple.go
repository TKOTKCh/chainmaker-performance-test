package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
	"io/ioutil"
)

const (
	projectInfoArgKey       = "projectInfo"
	projectIdArgKey         = "projectId"
	projectItemIdArgKey     = "itemId"
	projectInfoStoreKey     = "project"
	projectItemsStoreMapKey = "projectItems"
	projectVotesStoreMapKey = "projectVotes"
	trueString              = "1"
)

type ProjectInfo struct {
	Id        string      `json:"Id"`
	PicUrl    string      `json:"PicUrl"`
	Title     string      `json:"Title"`
	StartTime string      `json:"StartTime"`
	EndTime   string      `json:"EndTime"`
	Desc      string      `json:"Desc"`
	Items     []*ItemInfo `json:"Items"`
}

type ItemInfo struct {
	Id     string `json:"Id"`
	PicUrl string `json:"PicUrl"`
	Desc   string `json:"Desc"`
	Url    string `json:"Url"`
}
type ProjectInfoWrapper struct { // 新增外层结构体
	ProjectInfo *ProjectInfo `json:"projectInfo"` // 注意字段名与 JSON 的 key 对应
}
type ProjectVotesInfo struct {
	ProjectId string           `json:"ProjectId"`
	ItemVotes []*ItemVotesInfo `json:"ItemVotes"`
}

type ItemVotesInfo struct {
	ItemId     string   `json:"ItemId"`
	VotesCount int      `json:"VotesCount"`
	Voters     []string `json:"Voters"`
}

func main() {
	sha256.New()
	methodName := "sum"
	wasmBytes, err := ioutil.ReadFile("./sum/sum-go.wasm")
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
	start, err := instance.Exports.GetWasiStartFunction()
	if err != nil {
		fmt.Sprintf("Error getting WASI start function: %v", err)
	} else {
		start()
	}

	for i := 0; i < 10; i++ {
		// 获取并执行 WASI start 函数

		// 获取并执行 method 函数
		method, err := instance.Exports.GetFunction(methodName)

		result, err := method(1, 2)
		if err != nil {
			panic(fmt.Sprintf("Error calling %s function: %v", methodName, err))
		}

		fmt.Println(result)
	}

}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
