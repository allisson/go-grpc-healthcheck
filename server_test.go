package healthcheck

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type fakeChecker struct {
	calls int
}

func (f *fakeChecker) Calls() int {
	return f.calls
}

func (f *fakeChecker) Check() error {
	f.calls++
	return nil
}

type fakeFailChecker struct {
	calls int
}

func (f *fakeFailChecker) Calls() int {
	return f.calls
}

func (f *fakeFailChecker) Check() error {
	f.calls++
	return errors.New("Check Fail")
}

func TestServer(t *testing.T) {
	t.Run("AddChecker", func(t *testing.T) {
		expectedCheckers := map[string]Checker{
			"checker":      &fakeChecker{},
			"fail-checker": &fakeChecker{},
		}
		server := NewServer()
		server.AddChecker("checker", &fakeChecker{})
		server.AddChecker("fail-checker", &fakeChecker{})

		if !reflect.DeepEqual(expectedCheckers, server.checkers) {
			t.Errorf("got=%#v, want=%#v", server.checkers, expectedCheckers)
		}
	})

	t.Run("Check with empty checkers", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "unknown service")
		server := NewServer()
		request := healthpb.HealthCheckRequest{Service: ""}
		_, err := server.Check(context.Background(), &request)

		if !reflect.DeepEqual(expectedError, err) {
			t.Errorf("got=%#v, want=%#v", err, expectedError)
		}
	})

	t.Run("Check with nil checkers", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "unknown service")
		server := Server{}
		request := healthpb.HealthCheckRequest{Service: ""}
		_, err := server.Check(context.Background(), &request)

		if !reflect.DeepEqual(expectedError, err) {
			t.Errorf("got=%#v, want=%#v", err, expectedError)
		}
	})

	t.Run("Check with empty req.Service", func(t *testing.T) {
		expectedResponse := &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}
		server := NewServer()
		checker := &fakeChecker{}
		checkerTwo := &fakeChecker{}
		server.AddChecker("checker", checker)
		server.AddChecker("checker-two", checkerTwo)
		request := healthpb.HealthCheckRequest{Service: ""}
		response, err := server.Check(context.Background(), &request)

		if err != nil {
			t.Errorf("got=%#v, want=%#v", err, nil)
		}

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("got=%#v, want=%#v", response, expectedResponse)
		}

		if checker.Calls() != 1 {
			t.Errorf("got=%#v, want=%#v", checker.Calls(), 1)
		}

		if checkerTwo.Calls() != 1 {
			t.Errorf("got=%#v, want=%#v", checkerTwo.Calls(), 1)
		}
	})

	t.Run("Check with req.Service equals checker-two", func(t *testing.T) {
		expectedResponse := &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}
		server := NewServer()
		checker := &fakeChecker{}
		checkerTwo := &fakeChecker{}
		server.AddChecker("checker", checker)
		server.AddChecker("checker-two", checkerTwo)
		request := healthpb.HealthCheckRequest{Service: "checker-two"}
		response, err := server.Check(context.Background(), &request)

		if err != nil {
			t.Errorf("got=%#v, want=%#v", err, nil)
		}

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("got=%#v, want=%#v", response, expectedResponse)
		}

		if checker.Calls() != 0 {
			t.Errorf("got=%#v, want=%#v", checker.Calls(), 0)
		}

		if checkerTwo.Calls() != 1 {
			t.Errorf("got=%#v, want=%#v", checkerTwo.Calls(), 1)
		}
	})

	t.Run("Check with checker fail", func(t *testing.T) {
		expectedResponse := &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}
		server := NewServer()
		checker := &fakeChecker{}
		failChecker := &fakeFailChecker{}
		server.AddChecker("checker", checker)
		server.AddChecker("fail-checker", failChecker)
		request := healthpb.HealthCheckRequest{Service: "fail-checker"}
		response, err := server.Check(context.Background(), &request)

		if err != nil {
			t.Errorf("got=%#v, want=%#v", err, nil)
		}

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("got=%#v, want=%#v", response, expectedResponse)
		}

		if checker.Calls() != 0 {
			t.Errorf("got=%#v, want=%#v", checker.Calls(), 0)
		}

		if failChecker.Calls() != 1 {
			t.Errorf("got=%#v, want=%#v", failChecker.Calls(), 1)
		}
	})

	t.Run("Check with not found checker", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "unknown service")
		server := NewServer()
		checker := &fakeChecker{}
		server.AddChecker("checker", checker)
		request := healthpb.HealthCheckRequest{Service: "not-found-checker"}
		_, err := server.Check(context.Background(), &request)

		if !reflect.DeepEqual(expectedError, err) {
			t.Errorf("got=%#v, want=%#v", err, expectedError)
		}

		if checker.Calls() != 0 {
			t.Errorf("got=%#v, want=%#v", checker.Calls(), 0)
		}
	})

	t.Run("Watch", func(t *testing.T) {
		expectedError := status.Errorf(codes.Unimplemented, "method Watch not implemented")
		server := NewServer()
		request := healthpb.HealthCheckRequest{Service: ""}
		err := server.Watch(&request, nil)

		if !reflect.DeepEqual(expectedError, err) {
			t.Errorf("got=%#v, want=%#v", err, expectedError)
		}
	})
}
