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

func (r *Router) RegisterPicturesRoutes(controller *PicturesController) *Router {
	pictures := r.PathPrefix("/pictures").Subrouter()

	pictures.HandleFunc("", controller.pictures).Methods(http.MethodGet, http.MethodOptions).Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}")
	pictures.HandleFunc("", controller.pictureOfTheDay).Methods(http.MethodGet, http.MethodOptions).Queries("date", "{date:\\d{4}-\\d{2}-\\d{2}}")

	return r
}
