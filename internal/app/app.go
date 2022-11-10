package app

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/internal/config"
	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest"
)

func Init() {
	globalConfig := config.ParseConfig()

	_ = NewPostgresConnector().OpenDBConnect(globalConfig.DB())

	r, err := rest.NewRouter()
	if err != nil {
		log.Panicf("%+v", err)
	}

	//
	// Http server
	//

	server := http.Server{
		Addr:         globalConfig.App().Address(),
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Info("http server started at ", globalConfig.App().Address())
	}()

	log.Panicf("%+v", server.ListenAndServe())
}
