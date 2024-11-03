// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	logCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "log_count",
		Help: "The total number of logs received.",
	}, []string{"api_version", "log_type"})
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(logCount)
}

func IncLogCount(apiVersion, logType string, count float64) {
	logCount.WithLabelValues(apiVersion, logType).Add(count)
}
