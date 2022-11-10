package config

import "os"

type App struct {
	address string
}

type AppI interface {
	Address() string
}

func (app App) Address() string {
	return app.address
}

func (app *App) getEnv() error {
	app.address = os.Getenv(`APP_ADDRESS`)
	if app.address == "" {
		return getEnvErr(`APP_ADDRESS`)
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
