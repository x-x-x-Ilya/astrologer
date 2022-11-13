package config

import "os"

type App struct {
	address string
	apiKey  string
}

type AppI interface {
	Address() string
	APIKey() string
}

func (app App) Address() string {
	return app.address
}

func (app App) APIKey() string {
	return app.apiKey
}

func (app *App) getEnv() error {
	app.address = os.Getenv(`APP_ADDRESS`)
	if app.address == "" {
		return getEnvErr(`APP_ADDRESS`)
	}

	app.apiKey = os.Getenv(`API_KEY`)
	if app.apiKey == "" {
		return getEnvErr(`API_KEY`)
	}

	return nil
}

func newAppConfig() (AppI, error) {
	var app App

	err := app.getEnv()
	if err != nil {
		return nil, err
	}

	return app, nil
}
