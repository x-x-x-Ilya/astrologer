package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func NewRouter() (*Router, error) {
	router := mux.NewRouter()

	headersOK := handlers.AllowedHeaders([]string{
		"Accept",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
	})

	methodsOK := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodOptions,
	})

	credentialsOK := handlers.AllowCredentials()

	router.Use(handlers.CORS(
		headersOK,
		methodsOK,
		credentialsOK,
	))

	return &Router{
		router,
	}, nil
}

// GET pictures?limit=1&offset=0 (картинки дня всего альбома)
// GET pictures?date=2022-11-11 (картинка дня на указанный день)

func (r *Router) RegisterPicturesRoutes(controller *Pictures) *Router {
	pictures := r.PathPrefix("/pictures").Subrouter()

	pictures.HandleFunc("", controller.pictures).Methods(http.MethodGet, http.MethodOptions).Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}")
	pictures.HandleFunc("", controller.pictureOfTheDay).Methods(http.MethodGet, http.MethodOptions).Queries("date", "{date:YYYY-MM-DD}")

	return r
}
