package app

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/internal/config"
	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"github.com/x-x-x-Ilya/astrologer/internal/services"
	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest"
)

func InitServer() (http.Server, error) {
	globalConfig := config.ParseConfig()
	dbConf := globalConfig.DB()

	err := migrationsUp(dbConf.Address(), dbConf.Port(), dbConf.User(), dbConf.Password(), dbConf.Name())
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	db := NewPostgresConnector().OpenDBConnect(globalConfig.DB())

	transaction, err := database.NewTransaction(db)
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	picturesRepository, err := database.NewPicturesRepository(db)
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	transactionService, err := services.NewTransactionService(transaction)
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	clientService := services.NewClientService(time.Second * 10)

	nasaClient, err := services.NewNasaClient(globalConfig.App().APIKey(), clientService)
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	storageService := services.NewStorageService()

	picturesService, err := services.NewPicturesService(storageService, picturesRepository, nasaClient, transactionService)
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	picturesController, err := rest.NewPicturesController(picturesService)
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	r, err := rest.NewRouter()
	if err != nil {
		return http.Server{}, errors.WithStack(err)
	}

	r.RegisterPicturesRoutes(picturesController)

	return http.Server{
		Addr:         globalConfig.App().Address(),
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		Handler:      r,
	}, nil
}
