package main

import (
	"log"
	"net"

	alsv2 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	"github.com/zirain/als/pkg/Wasm/Common"
	"google.golang.org/grpc"
)

type ALSServer struct {
}

func (a *ALSServer) StreamAccessLogs(logStream alsv2.AccessLogService_StreamAccessLogsServer) error {
	log.Println("Streaming als v2 logs")
	for {
		data, err := logStream.Recv()
		if err != nil {
			return err
		}

		httpLogs := data.GetHttpLogs()
		if httpLogs != nil {
			for _, l := range httpLogs.LogEntry {
				upstream := l.CommonProperties.FilterStateObjects["wasm.upstream_peer"]
				if len(upstream.GetValue()) != 0 {
					node := Common.GetRootAsFlatNode(upstream.GetValue(), 3)
					log.Printf("wasm.upstream_peer workloadname: %s", node.WorkloadName())
				}
			}
		}

		//log.Printf("Received log data: %s\n", data.String())
	}
}

func NewALSServer() *ALSServer {
	return &ALSServer{}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to start listener on port 8080: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	alsv2.RegisterAccessLogServiceServer(grpcServer, NewALSServer())
	log.Println("Starting v2 ALS Server")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc serve err: %v", err)
	}
}
