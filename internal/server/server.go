package server

import (
	"QueryGuard/internal/adapter"
	"QueryGuard/internal/config"
	"QueryGuard/internal/detector"
	"QueryGuard/internal/extractor"
	"QueryGuard/internal/metrics"
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

func (s *Server) Start(){
	
}