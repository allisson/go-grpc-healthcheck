# go-grpc-healthcheck
[![Build Status](https://travis-ci.org/allisson/go-grpc-healthcheck.svg)](https://travis-ci.org/allisson/go-grpc-healthcheck) [![Go Report Card](https://goreportcard.com/badge/github.com/allisson/go-grpc-healthcheck)](https://goreportcard.com/report/github.com/allisson/go-grpc-healthcheck)

Simple implementation of `gRPC Health Checking Protocol`.

## How to use

```go
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
```

```bash
go run main.go
```

```bash
grpc_health_probe-linux-amd64 -addr=localhost:50051
service unhealthy (responded with "NOT_SERVING")
grpc_health_probe-linux-amd64 -addr=localhost:50051
status: SERVING
grpc_health_probe-linux-amd64 -addr=localhost:50051 -service "not-found-checker"
error: health rpc failed: rpc error: code = NotFound desc = unknown service
```

## References

- [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
- [grpc_health_probe](https://github.com/grpc-ecosystem/grpc-health-probe)
