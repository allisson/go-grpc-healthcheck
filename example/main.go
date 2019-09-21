package main

import (
	"errors"
	"log"
	"net"

	healthcheck "github.com/allisson/go-grpc-healthcheck"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// dumbChecker implements healthcheck.Checker interface
type dumbChecker struct {
	fail bool
}

func (d *dumbChecker) Check() error {
	d.fail = !d.fail
	if d.fail {
		return errors.New("Fail")
	}
	return nil
}

func main() {
	checker := dumbChecker{}
	healthcheckServer := healthcheck.NewServer()
	healthcheckServer.AddChecker("dumb-checker", &checker)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed-to-listen-tcp-port")
	}
	grpcServer := grpc.NewServer()
	healthpb.RegisterHealthServer(grpcServer, &healthcheckServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("failed-to-serve-grpc-requests")
	}
}
