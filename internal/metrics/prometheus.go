package metrics

import(
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct{
	CusumValue prometheus.Gauge
	EwmaValue prometheus.Gauge
	CusumAlertTotal prometheus.Counter
    EwmaAlertTotal prometheus.Counter
	DriftAlertTotal prometheus.Counter
	FalsePositiveAlert prometheus.Gauge
	RequestLatency prometheus.Histogram
	RequestsTotal prometheus.Counter
	RagForwardTotal prometheus.Counter
	KbUpdateTotal      prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *Metrics{
	m := &Metrics{
		CusumValue: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:"cusum_values",
			Help: "Current value of Cusum",
		}),
		EwmaValue: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:"ewma_values",
			Help: "Current value of Ewma",
		}),
		FalsePositiveAlert: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:"false_positive_alert",
			Help: "Alert on stable baseline",
		}),
		CusumAlertTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cusum_alert_total",
			Help:"Total number of time cusum alert happened",
		}),
		EwmaAlertTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ewma_alert_total",
			Help:"Total number of time ewma alert happened",
		}),
		DriftAlertTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "drift_alert_total",
			Help:"Total number of drift events happened",
		}),
		RequestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "request_total",
			Help:"Total number of request",
		}),
		RagForwardTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "rag_forward_total",
			Help:"Total number of time query forwaded to rag/model",
		}),
		KbUpdateTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kb_update_total",
			Help:"Total number of knowledge based updates",
		}),
		RequestLatency: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "queryguard_request_latency_ms",
            Help:    "Latency per request in ms",
            Buckets: []float64{10, 50, 100, 200, 500},
        }),
	}

	reg.MustRegister(
		m.CusumValue,
		m.EwmaValue,
		m.CusumAlertTotal,
		m.EwmaAlertTotal,
		m.DriftAlertTotal,
		m.FalsePositiveAlert,
		m.RequestLatency,
		m.RequestsTotal,
		m.RagForwardTotal,
		m.KbUpdateTotal,
	)
	return m
}


func (m *Metrics) RecordDrift(cusumScore, ewmaScore float64, drifted bool) {
    m.CusumValue.Set(cusumScore)
    m.EwmaValue.Set(ewmaScore)
    if drifted {
        m.DriftAlertTotal.Inc()
    }
}

func (m *Metrics) RecordLatency(ms float64) {
    m.RequestLatency.Observe(ms)
}

func (m *Metrics) RecordRequest() {
    m.RequestsTotal.Inc()
}