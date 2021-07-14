#!/bin/bash

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./docs/
cd web && GOOS=js GOARCH=wasm go build -o  ../docs/sim.wasm