package rest_server

import (
	"fmt"
	"net/http"

	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal/rest_server/middleware"
	"github.com/gorilla/mux"
)

type RESTServer struct {
	conf   *config.Config
	server *http.Server
}

func NewServer(conf *config.Config) internal.ServerImpl {
	return &RESTServer{
		conf:   conf,
		server: &http.Server{Addr: ":" + fmt.Sprint(conf.RESTPort)},
	}
}

// Serve registers the services and starts serving
func (s *RESTServer) Serve() error {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// register handlers
	router.HandleFunc("/keys", s.keyExchangeHandler).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/data", s.DataTransferHandler).Methods(http.MethodPost, http.MethodOptions)

	// start server
	router.Use(middleware.Cors)
	s.server.Handler = router
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *RESTServer) Shutdown() {

}
