package middleware_test

import (
	"go-vue-journey/internal/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCORS_SetsAccessControlAllowOriginHeader(t *testing.T) {
	tests := []struct {
		name        string
		frontendURL string
	}{
		{
			name:        "sets header from FrontendURL config",
			frontendURL: "http://localhost:3000",
		},
		{
			name:        "sets header with different origin",
			frontendURL: "https://example.com",
		},
		{
			name:        "sets header with wildcard",
			frontendURL: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("FRONTEND_URL", tt.frontendURL)
			defer os.Unsetenv("FRONTEND_URL")

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.CORS(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			gotOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if gotOrigin != tt.frontendURL {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", gotOrigin, tt.frontendURL)
			}

			gotMethods := w.Header().Get("Access-Control-Allow-Methods")
			expectedMethods := "GET, POST, PUT, DELETE, OPTIONS"
			if gotMethods != expectedMethods {
				t.Errorf("Access-Control-Allow-Methods = %q, want %q", gotMethods, expectedMethods)
			}

			gotHeaders := w.Header().Get("Access-Control-Allow-Headers")
			expectedHeaders := "Content-Type, Authorization"
			if gotHeaders != expectedHeaders {
				t.Errorf("Access-Control-Allow-Headers = %q, want %q", gotHeaders, expectedHeaders)
			}

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestCORS_OptionsRequestReturnsNoContent(t *testing.T) {
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	defer os.Unsetenv("FRONTEND_URL")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for OPTIONS request")
	})

	handler := middleware.CORS(nextHandler)

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", w.Code)
	}

	gotOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if gotOrigin != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", gotOrigin, "http://localhost:3000")
	}
}

func TestCORS_NonOptionsRequestCallsNextHandler(t *testing.T) {
	os.Setenv("FRONTEND_URL", "http://localhost:8080")
	defer os.Unsetenv("FRONTEND_URL")

	nextHandlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("response body"))
	})

	handler := middleware.CORS(nextHandler)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			nextHandlerCalled = false

			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if !nextHandlerCalled {
				t.Error("next handler should be called for non-OPTIONS request")
			}

			if w.Code != http.StatusCreated {
				t.Errorf("expected status 201, got %d", w.Code)
			}

			if w.Body.String() != "response body" {
				t.Errorf("expected body 'response body', got %q", w.Body.String())
			}

			gotOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if gotOrigin != "http://localhost:8080" {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", gotOrigin, "http://localhost:8080")
			}
		})
	}
}

func TestCORS_EmptyFrontendURL(t *testing.T) {
	os.Unsetenv("FRONTEND_URL")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.CORS(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	gotOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if gotOrigin != "" {
		t.Errorf("Access-Control-Allow-Origin = %q, want empty string", gotOrigin)
	}
}
