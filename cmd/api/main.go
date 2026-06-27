package main

import (
	"log"
    "QueryGuard/internal/adapter"
    "QueryGuard/internal/config"
    "QueryGuard/internal/detector"
    "QueryGuard/internal/extractor"
    "QueryGuard/internal/metrics"
    "QueryGuard/internal/server"
    "github.com/prometheus/client_golang/prometheus"
)   



func main(){
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    ext, err := extractor.NewExtractor(cfg.Baseline.Store)
    if err != nil {
        log.Fatal(err)
    }

    det := &detector.Detector{
        Cusum: detector.CusumDetector{
        Target:    float64(cfg.Metrics.CUSUM.Mean),
        Slack:     float64(cfg.Metrics.CUSUM.K),
        Threshold: float64(cfg.Metrics.CUSUM.H),
    },
    Ewma: detector.EwmaDetector{
        Lambda:    float64(cfg.Metrics.EWMA.Lambda),
        Current:   float64(cfg.Metrics.EWMA.Mean),
        Mean:      float64(cfg.Metrics.EWMA.Mean),
        StdDev:    float64(cfg.Metrics.EWMA.StdDev),
        Threshold: float64(cfg.Metrics.EWMA.Threshold),
    },
    }
    reg := prometheus.NewRegistry()
    m := metrics.NewMetrics(reg)
    a := adapter.NewAdapter(cfg)

    s := server.New(cfg, ext, det, m, a,reg)
    s.Start()
}