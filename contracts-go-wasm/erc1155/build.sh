#!/bin/bash

contractName=$1
buildOption=$2
targetARCH=$3
crypto=""

if [ "$(uname)" == "Linux" ];then
  crypto="-tags crypto"
fi

if  [[ ! -n $contractName ]] ;then
    echo "contractName is empty. use as: ./build.sh contractName."
    exit 1
fi

# 默认 buildOption 为 tinygo
if [[ -z $buildOption ]]; then
    buildOption="tinygo"
fi

# 根据 buildOption 选择编译方式
if [[ $buildOption == "tinygo" ]]; then

    export GOROOT="/home/chenhang/WorkSpace/go-v1.23.7"
    export PATH="$GOROOT/bin:$PATH"
    echo "Using TinyGo to compile..."
#    tinygo build -no-debug -opt=s -o "$contractName-tinygo3.wasm" -target wasi
    GOOS=wasm GOARCH=wasip1 tinygo build -buildmode=c-shared -o "$contractName-tinygo.wasm" -target wasi
else
#  要用go自带的编译器编译成wasm要切换成1.24.1版本
    export GOROOT="/home/chenhang/WorkSpace/go-v1.24.1"
    export PATH="$GOROOT/bin:$PATH"
    echo "Using Go to compile..."
    GOOS=wasip1 GOARCH=wasm go build -o "$contractName-go.wasm"
fi





