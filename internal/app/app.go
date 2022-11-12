package app

import (
	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"github.com/x-x-x-Ilya/astrologer/internal/services"
	"net/http"
	"time"

	"github.com/x-x-x-Ilya/astrologer/internal/config"
	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest"
)

func InitServer() (http.Server, error) {
	globalConfig := config.ParseConfig()
	dbConf := globalConfig.DB()

	err := migrationsUp(dbConf.Address(), dbConf.Port(), dbConf.User(), dbConf.Password(), dbConf.Name())
	if err != nil {
		return http.Server{}, err
	}

	db := NewPostgresConnector().OpenDBConnect(globalConfig.DB())

	picturesRepository, err := database.NewPicturesRepository(db)
	if err != nil {
		return http.Server{}, err
	}

	clientService := services.NewClientService(time.Second * 10)

	nasaClient, err := services.NewNasaClient(globalConfig.App().ApiKey(), clientService)
	if err != nil {
		return http.Server{}, err
	}

	storageService := services.NewStorageService()

	picturesService, err := services.NewPicturesService(storageService, picturesRepository, nasaClient)
	if err != nil {
		return http.Server{}, err
	}

	picturesController, err := rest.NewPicturesController(picturesService)
	if err != nil {
		return http.Server{}, err
	}

	r, err := rest.NewRouter()
	if err != nil {
		return http.Server{}, err
	}

	r.RegisterPicturesRoutes(picturesController)

	return http.Server{
		Addr:         globalConfig.App().Address(),
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		Handler:      r,
	}, nil
}
