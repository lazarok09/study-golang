package router

import (
	"api/src/routes"

	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	return routes.ConfigRoutes(r)
}
