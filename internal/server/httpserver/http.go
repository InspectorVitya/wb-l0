package httpserver

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/inspectorvitya/wb-l0/internal/application"
	"github.com/inspectorvitya/wb-l0/internal/config"
	"net"
	"net/http"
	_ "net/http/pprof"
)

type Server struct {
	App        *application.App
	router     *mux.Router
	httpServer *http.Server
}

func New(cfg config.HTTP, App *application.App) *Server {
	router := mux.NewRouter()
	server := &Server{
		httpServer: &http.Server{
			Addr:    net.JoinHostPort("", cfg.Port),
			Handler: router,
		},
		router: router,
		App:    App,
	}

	return server
}

func (s *Server) Start() error {
	vueHandler := http.FileServer(http.Dir("./web/dist"))
	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	s.router.HandleFunc("/order/{id}", s.GetOrder).Methods(http.MethodGet)
	s.router.PathPrefix("/").Handler(vueHandler)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
