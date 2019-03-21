package main

import (
	"fmt"

	"github.com/ymgyt/happycode/frontend"
)

func main() {
	fmt.Println("start app")

	app := frontend.NewApp()
	app.Init()
	fmt.Println("init app")

	select {}
}
