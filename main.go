package main

import (
	"bookateria-api-go/account"
	"bookateria-api-go/auth"
	"bookateria-api-go/document"
	"bookateria-api-go/forum"
	"bookateria-api-go/log"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	// Documentation route
	fs := http.FileServer(http.Dir("./docs"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	// Routes for Documents
	versionRouter := router.PathPrefix("/v1").Subrouter()
	subRouter := versionRouter.PathPrefix("/document").Subrouter()
	subRouter.HandleFunc("", document.GetDocuments).Methods("GET")
	subRouter.HandleFunc("/{id}", document.GetDocument).Methods("GET")
	subRouter.HandleFunc("", document.PostDocument).Methods("POST")
	subRouter.HandleFunc("/{id}", document.UpdateDocument).Methods("POST")
	subRouter.HandleFunc("/{id}", document.DeleteDocument).Methods("DELETE")

	account.Router(versionRouter.PathPrefix("/account").Subrouter())
	auth.Router(versionRouter.PathPrefix("/auth").Subrouter())
	forum.Router(versionRouter.PathPrefix("/forum").Subrouter())

	loggerMgr := log.InitLog()
	zap.ReplaceGlobals(loggerMgr)
	logger := loggerMgr.Sugar()
	logger.Debug("Starting Server")
	logger.Debug(http.ListenAndServe(":5000", router))
	err := loggerMgr.Sync()
	logger.Debug("Buffer flush issue: ", err)
}
