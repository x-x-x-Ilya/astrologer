package app

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/config"
	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"github.com/x-x-x-Ilya/astrologer/internal/services"
	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest"
)

const (
	readTimeout  = time.Second * 60
	writeTimeout = time.Second * 60
)

func InitServer() (*http.Server, error) {
	var (
		globalConfig = config.ParseConfig()
		dbConfig     = globalConfig.DB()
		appConfig    = globalConfig.App()
	)

	err := migrationsUp(dbConfig.ConnectionString())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db := NewPostgresConnector().OpenDBConnect(dbConfig)

	transaction, err := database.NewTransaction(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	picturesRepository, err := database.NewPicturesRepository(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	transactionService, err := services.NewTransactionService(transaction)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clientService := services.NewClientService(time.Second * 10)

	nasaClient, err := services.NewNasaClient(appConfig.APIKey(), clientService)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	storageService := services.NewStorageService(appConfig.Storage())

	picturesService, err := services.NewPicturesService(storageService, picturesRepository, nasaClient, transactionService)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	picturesController, err := rest.NewPicturesController(picturesService)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r, err := rest.NewRouter()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.RegisterPicturesRoutes(picturesController)

	return &http.Server{
		Addr:         appConfig.Address(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      r,
	}, nil
}
