package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/x-x-x-Ilya/astrologer/internal/app"
	_ "github.com/x-x-x-Ilya/astrologer/internal/app"
)

func main() {
	server, err := app.InitServer()
	if err != nil {
		log.Panicf("%+v", err.Error())
	}

	log.Info("Server has been started")

	log.Panicf("%+v", server.ListenAndServe())
}
