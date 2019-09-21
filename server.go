package healthcheck

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Server is a implementation of healthpb.HealthServer
type Server struct {
	checkers map[string]Checker
}

// AddChecker includes a new checker on the checkers maps
func (s *Server) AddChecker(name string, checker Checker) {
	s.checkers[name] = checker
}

func (s *Server) checkersIsEmpty() bool {
	if s.checkers == nil {
		return true
	} else if len(s.checkers) == 0 {
		return true
	}
	return false
}

// Check is a implementation of healthpb.HealthServer.Check
func (s *Server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if s.checkersIsEmpty() {
		log.Println("healthcheck-checkers-is-empty")
		return nil, status.Error(codes.NotFound, "unknown service")
	}

	var servicesToCheck []string

	if req.Service == "" {
		for name := range s.checkers {
			servicesToCheck = append(servicesToCheck, name)
		}
	} else {
		servicesToCheck = []string{req.Service}
	}

	for _, serviceName := range servicesToCheck {
		checker, ok := s.checkers[serviceName]
		if !ok {
			log.Printf("healthcheck-checker-not-found, name=%s\n", serviceName)
			return nil, status.Error(codes.NotFound, "unknown service")
		}
		if err := checker.Check(); err != nil {
			log.Printf("healthcheck-checker-error, name=%s, error=%s\n", serviceName, err.Error())
			return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}, nil
		}
	}

	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

// Watch is a implementation of healthpb.HealthServer.Watch
func (s *Server) Watch(req *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

// NewServer returns a Server with initialized checkers maps
func NewServer() Server {
	return Server{checkers: make(map[string]Checker)}
}
