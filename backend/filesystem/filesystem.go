package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/core/errors"
	"github.com/ymgyt/happycode/core/errors/errcode"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*config.Config, error) {
	log.V(0).Debug("load config", zap.String("path", path))
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WrapM(errcode.FileNotFound, err, path)
		}
		// TODO permission error handling
		panic(err)
	}

	return config.NewFromYAML(b, path)
}

func SaveConfig(cfg *config.Config, path string) error {
	// TODO permission error handling
	if err := os.MkdirAll(filepath.Dir(path), cfg.DirPermission()); err != nil {
		panic(err)
	}
	b, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, b, cfg.FilePermission())
	if err != nil {
		panic(err)
	}
	return nil
}
