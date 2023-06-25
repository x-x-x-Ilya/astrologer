package config

import "os"

type app struct {
	address string
	apiKey  string
	storage string
}

type AppI interface {
	Address() string
	APIKey() string
	Storage() string
}

func (app app) Address() string {
	return app.address
}

func (app app) APIKey() string {
	return app.apiKey
}

func (app app) Storage() string {
	return app.storage
}

func (app *app) fillFromEnv() error {
	const (
		Address = `APP_ADDRESS`
		APIKey  = `API_KEY`
		Storage = `APP_STORAGE`
	)

	app.address = os.Getenv(Address)
	if app.address == "" {
		return getEnvErr(Address)
	}

	app.apiKey = os.Getenv(APIKey)
	if app.apiKey == "" {
		return getEnvErr(APIKey)
	}

	app.storage = os.Getenv(Storage)
	if app.apiKey == "" {
		return getEnvErr(Storage)
	}

	_ = os.Mkdir(app.storage, os.ModePerm) // creates storage folder if it not exists

	return nil
}

func newAppConfig() (AppI, error) {
	var app app

	err := app.fillFromEnv()
	if err != nil {
		return nil, err
	}

	return app, nil
}
