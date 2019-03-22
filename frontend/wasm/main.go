package main

import (
	"go.uber.org/zap"

	"github.com/ymgyt/happycode/core"
	"github.com/ymgyt/happycode/frontend/di"
	"github.com/ymgyt/happycode/frontend/log"
)

func main() {
	log.V(0).Info("start app", zap.String("version", core.Version))
	di.NewApp().Run()
}
