package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewerPageGET(t *testing.T) {
	// Load templates
	if err := loadTemplates(); err != nil {
		t.Errorf("Failed to load templates: %v", err)
	}

	// Test GET request
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(viewerHandler).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Viewer page returned %v != %v on GET", w.Code, http.StatusOK)
	}

	// Test GET request on not found page
	r, _ = http.NewRequest("GET", "", nil)
	w = httptest.NewRecorder()
	http.HandlerFunc(viewerHandler).ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Viewer page returned %v != %v on GET", w.Code, http.StatusOK)
	}

	// Test GET request on statistics page
	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	http.HandlerFunc(statisticsHandler).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Viewer page returned %v != %v on GET", w.Code, http.StatusOK)
	}
}
