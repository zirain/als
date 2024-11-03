// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	alsv2 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	alsv3 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"

	v2 "github.com/zirain/als/alsserver/v2"
	v3 "github.com/zirain/als/alsserver/v3"
)

var (
	alsAddr          = flag.String("addr", ":8080", "gRPC port for the envoy ALS server")
	monitoringAddr   = flag.String("monitoringAddr", ":19001", "port for the monitoring server")
	otelExporterAddr = flag.String("otelExportAddr", "otel-collector.monitoring:4317", "address for the OpenTelemetry collector")
)

func main() {
	// Set up monitoring server
	mux := http.NewServeMux()
	if err := addMonitor(mux); err != nil {
		log.Printf("could not establish self-monitoring: %v\n", err)
	}
	s := &http.Server{
		Addr:    *monitoringAddr,
		Handler: mux,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("monitoring server failed: %v", err)
		}
	}()

	// Set up ALS gRPC Server
	listener, err := net.Listen("tcp", *alsAddr)
	if err != nil {
		log.Fatalf("Failed to start listener on port 8080: %v", err)
	}
	exp, err := otlploggrpc.New(context.Background(),
		otlploggrpc.WithInsecure(), otlploggrpc.WithEndpoint(*otelExporterAddr),
	)
	processor := sdklog.NewBatchProcessor(exp)
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}()
	global.SetLoggerProvider(provider)

	logger := provider.Logger("envoy-als")

	if err != nil {
		log.Fatalf("Failed to create OTLP log exporter: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	alsv2.RegisterAccessLogServiceServer(grpcServer, v2.New(logger))
	alsv3.RegisterAccessLogServiceServer(grpcServer, v3.New(logger))
	log.Printf("ALS Server receive logs from [%s] send to %s \n", *alsAddr, *otelExporterAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc serve err: %v", err)
	}
}

func addMonitor(mux *http.ServeMux) error {
	mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{EnableOpenMetrics: true}))
	return nil
}
