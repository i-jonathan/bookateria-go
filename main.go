package main

import (
	"bookateria-api-go/account"
	"bookateria-api-go/auth"
	"bookateria-api-go/document"
	"bookateria-api-go/log"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main()  {
	router := mux.NewRouter()	
	document.Router(router.PathPrefix("/document").Subrouter())
	account.Router(router.PathPrefix("/account").Subrouter())
	auth.Router(router.PathPrefix("/auth").Subrouter())

	loggerMgr := log.InitLog()
	zap.ReplaceGlobals(loggerMgr)
	logger := loggerMgr.Sugar()
	logger.Debug("Starting Server")
	logger.Debug(http.ListenAndServe(":5000", router))
	err := loggerMgr.Sync()
	logger.Debug("Buffer flush issue: ", err)
}
