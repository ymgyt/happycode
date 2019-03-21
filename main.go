package main

import (
	"context"

	"github.com/ymgyt/happycode/cli"
)

func main() {
	cli.New().Execute(context.Background())
}
