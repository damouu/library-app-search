package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeDependencyChecker struct {
	err error
}

func (f *fakeDependencyChecker) Check() error {
	return f.err
}

func TestLiveHandler_ReturnsOK(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	recorder := httptest.NewRecorder()

	LiveHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestReadyHandler_ReturnsOKWhenDependencyIsHealthy(t *testing.T) {
	checker := &fakeDependencyChecker{}
	handler := NewHandler(checker)

	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	recorder := httptest.NewRecorder()

	handler.ReadyHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestReadyHandler_ReturnsServiceUnavailableWhenDependencyFails(t *testing.T) {
	checker := &fakeDependencyChecker{
		err: errors.New("elasticsearch unavailable"),
	}

	handler := NewHandler(checker)

	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	recorder := httptest.NewRecorder()

	handler.ReadyHandler(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			recorder.Code,
		)
	}
}
