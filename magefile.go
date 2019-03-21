// +build mage

package main

import "github.com/magefile/mage/sh"

// var Aliases = map[string]interface{}{}

func All() {
	Wasm()
	Run()
}

// Wasm build webassembly.
func Wasm() {
	sh.RunWith(map[string]string{
		"GOOS":   "js",
		"GOARCH": "wasm",
	}, "go", "build", "-o", "static/wasm/main.wasm", "frontend/wasm/main.go")
}

// Run run server.
func Run() {
	sh.RunV("go", "run", "main.go")
}

// Tidy prune any no-longer needed dependencies from go.mod
func Tidy() {
	sh.RunV("go", "mod", "tidy")
}

