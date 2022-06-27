
// main.go
package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"api/services"
)

type App struct {
	Router *mux.Router
	ArticleService *services.ArticleService
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter().StrictSlash(true)

	a.ArticleService = &services.ArticleService{
		Router: a.Router,
	}
	a.ArticleService.Initialize()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}