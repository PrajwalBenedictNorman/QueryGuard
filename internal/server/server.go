package server

import (
	"QueryGuard/internal/adapter"
	"QueryGuard/internal/config"
	"QueryGuard/internal/detector"
	"QueryGuard/internal/extractor"
	"QueryGuard/internal/metrics"
	"fmt"
	"log"
	"net/http"
)

type Server struct{
	Config *config.Config
	Extractor *extractor.Extactor
	Detector  *detector.Detector
    Metrics   *metrics.Metrics
    Adapter   *adapter.Adapter
}

func New(
    cfg *config.Config,
    ext *extractor.Extactor,
    det *detector.Detector,
    m *metrics.Metrics,
    a *adapter.Adapter,
) *Server{
	return &Server{
		Config: cfg,
		Extractor: ext,
		Detector: det,
		Metrics: m,
		Adapter: a,
	}
}

func (s *Server) Start() {
	s.registerRoutes()
	port := fmt.Sprintf(":%d", s.Config.Server.Port)
	log.Printf("QueryGuard started on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}