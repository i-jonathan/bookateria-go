package main

import (
	"bookateriago/account"
	"bookateriago/assignment"
	"bookateriago/auth"
	"bookateriago/document"
	"bookateriago/forum"
	"bookateriago/log"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
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

	router.Use(corsMiddleware)

	log.Start("Starting Server")
	err := http.ListenAndServe(":5000", router)
	log.ErrorHandler(err)

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://bookateria.net")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
