package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/Hitomiblood/StockStream/internal/handlers"
	"github.com/gin-gonic/gin"
)

func TestSetupRouter_RootRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter(&config.Config{LogLevel: "debug"}, &handlers.StockHandler{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}
	if location := w.Header().Get("Location"); location != "/swagger/index.html" {
		t.Fatalf("location = %q", location)
	}
}

func TestSetupRouter_RegistersExpectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter(&config.Config{LogLevel: "debug"}, &handlers.StockHandler{})

	routes := r.Routes()
	if len(routes) < 10 {
		t.Fatalf("expected many routes registered, got %d", len(routes))
	}
}
