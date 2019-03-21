package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ymgyt/cli"
	"github.com/ymgyt/happycode/backend/di"
	"github.com/ymgyt/happycode/backend/filesystem"
	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/core/errors"
	"github.com/ymgyt/happycode/core/errors/errcode"
	"go.uber.org/zap"
)

func New() *cli.Command {
	root := &cli.Command{
		Name: "happycode",
		Run:  Run,
	}

	return root
}

func Run(ctx context.Context, _ *cli.Command, _ []string) {
	ctx, cancel := context.WithCancel(context.Background())
	go watchSignal(cancel)

	app := di.NewAppStack(loadConfig())
	app.Server.Run(ctx)
}

func loadConfig() *config.Config {
	path := config.DefaultConfigPath()
	cfg, err := filesystem.LoadConfig(path)
	if errors.Code(err) == errcode.FileNotFound {
		// initial flow, setup happycode.
		logger := log.V(0).With(zap.String("path", path))

		logger.Info("config file not found")
		logger.Info("create config file")
		cfg = config.Default()
		err = filesystem.SaveConfig(cfg, path)
		if err != nil {
			logger.Fatal("failed to create config file")
		}
		cfg.Meta.FilePath = path
	}
	return cfg
}

func watchSignal(cancel func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP)
	switch sig := <-ch; sig {
	default:
		log.V(5).Info("receive signal", zap.String("sig", sig.String()))
		cancel()
	}
}
