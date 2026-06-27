package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Answer  string  `json:"answer"`
	Drifted bool    `json:"drifted"`
	Score   float64 `json:"score"`
}

func (s *Server) registerRoutes() {
	http.HandleFunc("/query", s.handleQuery)
	http.Handle("/metrics", promhttp.Handler())
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	s.Metrics.RecordRequest()

	start := time.Now()

	score := s.Extractor.Extract(req.Query)

	result := s.Detector.Detection(score)

	s.Metrics.RecordDrift(result.CusumValue, result.EwmaValue, result.DriftAlert)
	if result.CusumAlert {
		s.Metrics.CusumAlertTotal.Inc()
	}
	if result.EwmaAlert {
		s.Metrics.EwmaAlertTotal.Inc()
	}

	answer, err := s.Adapter.Forward(req.Query)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusInternalServerError)
		return
	}
	s.Metrics.RagForwardTotal.Inc()

	s.Metrics.RecordLatency(float64(time.Since(start).Milliseconds()))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(QueryResponse{
		Answer:  answer,
		Drifted: result.DriftAlert,
		Score:   score,
	})
}