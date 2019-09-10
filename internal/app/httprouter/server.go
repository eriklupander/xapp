package httprouter

import (
	"github.com/callistaenterprise/xapp/cmd"
	"github.com/callistaenterprise/xapp/internal/app/persistence"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	cfg     cmd.HTTPServerConfig
	mux     chi.Router
	storage persistence.Database
}

func NewServer(cfg *cmd.Config, storage persistence.Database) *Server {
	return &Server{cfg.HTTPServerConfig, chi.NewRouter(), storage}
}

// SetupRoutes initializes middlewares and declares routes.
// Note that this is a separate method in order to simplify unit-testing of HTTP routes and middlewares.
func (s *Server) SetupRoutes() {

	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.RealIP)
	//s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
	s.mux.Use(middleware.Timeout(time.Minute))

	s.mux.Get("/health", s.healthCheck)
}

func (s *Server) healthCheck(rw http.ResponseWriter, r *http.Request) {
	if err := s.storage.Ping(); err != nil {
		logrus.Warnf("Healtcheck failed performing DB check. Error: %+v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (s *Server) Start() {
	logrus.Infof("Starting HTTP server on %v", s.cfg.BindAddress)
	err := http.ListenAndServe(s.cfg.BindAddress, s.mux)
	if err != nil {
		logrus.WithError(err).Fatalf("error starting HTTP server on %v", s.cfg.BindAddress)
	}
}
