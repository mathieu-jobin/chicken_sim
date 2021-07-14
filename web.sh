#!/bin/bash

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./js/
cd web && GOOS=js GOARCH=wasm go build -o  ../js/sim.wasm