// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v3

import (
	"context"

	"log"

	alsv3 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	otellog "go.opentelemetry.io/otel/log"

	"github.com/zirain/als/monitoring"
)

type ALSServer struct {
	logger otellog.Logger
}

func New(logger otellog.Logger) *ALSServer {
	return &ALSServer{
		logger: logger,
	}
}

func (a *ALSServer) StreamAccessLogs(logStream alsv3.AccessLogService_StreamAccessLogsServer) error {
	log.Println("Streaming als v3 logs")
	for {
		data, err := logStream.Recv()
		if err != nil {
			return err
		}

		httpLogs := data.GetHttpLogs()
		if httpLogs != nil {
			monitoring.IncLogCount("v3", "http", float64(len(httpLogs.LogEntry)))
		}
		tcpLogs := data.GetTcpLogs()
		if tcpLogs != nil {
			monitoring.IncLogCount("v3", "tcp", float64(len(tcpLogs.LogEntry)))
		}

		if a.logger != nil {
			for _, record := range toLogRecord(data) {
				a.logger.Emit(context.TODO(), record)
			}
		}

		log.Printf("Received v3 log data: %s\n", data.String())
	}
}

var apiVersionAttr = otellog.String("api_version", "v3")

func toLogRecord(data *alsv3.StreamAccessLogsMessage) []otellog.Record {
	records := make([]otellog.Record, 0)
	httpLogs := data.GetHttpLogs()
	if httpLogs != nil {
		for _, httpLog := range httpLogs.LogEntry {
			var r otellog.Record
			r.AddAttributes(otellog.String("log_type", "http"), apiVersionAttr)
			r.SetBody(otellog.StringValue(httpLog.String()))
			records = append(records, r)
		}
	}

	tcpLogs := data.GetTcpLogs()
	if tcpLogs != nil {
		for _, tcpLog := range tcpLogs.LogEntry {
			var r otellog.Record
			r.AddAttributes(otellog.String("log_type", "tcp"), apiVersionAttr)
			r.SetBody(otellog.StringValue(tcpLog.String()))
			records = append(records, r)
		}
	}
	return records
}
