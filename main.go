package main

import (
	"bookateria-api-go/account"
	"bookateria-api-go/assignment"
	"bookateria-api-go/auth"
	"bookateria-api-go/document"
	"bookateria-api-go/forum"
	"bookateria-api-go/log"
	"github.com/gorilla/mux"
	"net/http"
)

func main()  {
	router := mux.NewRouter()
	// Documentation route
	fs := http.FileServer(http.Dir("./docs"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))
	versionRouter := router.PathPrefix("/v1").Subrouter()

	document.Router(versionRouter.PathPrefix("/document").Subrouter())
	account.Router(versionRouter.PathPrefix("/account").Subrouter())
	auth.Router(versionRouter.PathPrefix("/auth").Subrouter())
	forum.Router(versionRouter.PathPrefix("/forum").Subrouter())
	assignment.Router(versionRouter.PathPrefix("/assignment").Subrouter())

	log.AccessHandler("Starting server")
	err := http.ListenAndServe(":5000", router)
	log.ErrorHandler(err)

}
