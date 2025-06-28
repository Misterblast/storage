package utils

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Registry = prometheus.NewRegistry()

	UploadRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_upload_requests_total",
			Help: "Total number of file upload requests",
		},
		[]string{"destination", "status"},
	)

	UploadDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "file_upload_duration_seconds",
			Help:    "Histogram of file upload durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"destination"},
	)

	UploadFileSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "file_upload_size_bytes",
			Help:    "Histogram of uploaded file sizes",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 10),
		},
		[]string{"destination"},
	)

	UploadErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_upload_errors_total",
			Help: "Total number of file upload errors",
		},
		[]string{"destination", "error_type"},
	)
)

func Init() {
	Registry.MustRegister(UploadRequestCounter)
	Registry.MustRegister(UploadDuration)
	Registry.MustRegister(UploadFileSize)
	Registry.MustRegister(UploadErrorCounter)
}
