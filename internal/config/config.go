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

type Config struct {
	app AppI
	db  DBI
}

func newConfig() (ConfigI, error) {
	var (
		config Config
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

func (conf Config) App() AppI {
	return conf.app
}

func (conf Config) DB() DBI {
	return conf.db
}

func getEnvErr(key string) error {
	return errors.Errorf("env key %s not found", key)
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
