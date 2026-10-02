package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRecoveryMenangkapPanic: handler yang panic harus dibalas 500,
// dan program tidak boleh crash.
func TestRecoveryMenangkapPanic(t *testing.T) {
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("ada error tak terduga")
	})

	app := Recovery(panicking)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusInternalServerError)
	}
}

// TestLogging: middleware Logging tidak boleh mengubah response handler asli.
func TestLogging(t *testing.T) {
	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app := Logging(dummy)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusOK)
	}
}
