// Package config consist of files that contain config components(initialization and reading) interfaces with implementation
// and config.go file that combines all of them in one ConfigI that only available outside the package.
package config

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ConfigI interface {
	App() AppI
	DB() DBI
}

type config struct {
	app AppI
	db  DBI
}

func newConfig() (ConfigI, error) {
	var (
		config config
		err    error
	)

	config.app, err = newAppConfig()
	if err != nil {
		return config, errors.Wrap(err, "app env not found")
	}

	config.db, err = newDBConfig()
	if err != nil {
		return config, errors.Wrap(err, "db env not found")
	}

	return config, nil
}

func (conf config) App() AppI {
	return conf.app
}

func (conf config) DB() DBI {
	return conf.db
}

func ParseConfig() ConfigI {
	cp := flag.String("config", "./.env", "pass path to the config")

	flag.Parse()

	conf, err := loadConfig(*cp)
	if err != nil {
		log.Panicf("%+v", err)
	}

	return conf
}

func loadConfig(confPath string) (ConfigI, error) {
	_ = godotenv.Load(confPath)

	return newConfig()
}

func getEnvErr(key string) error {
	return errors.Errorf("env key %s not found", key)
}
