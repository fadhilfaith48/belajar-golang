package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestCORSHeader: setiap response harus punya header CORS, dan
// handler di akhirnya tetap dipanggil.
func TestCORSHeader(t *testing.T) {
	// bool ini membuktikan handler benar-benar dipanggil atau tidak
	handlerDipanggil := false
	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerDipanggil = true
		w.WriteHeader(http.StatusOK)
	})

	app := CORS(dummy)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	izin := w.Header().Get("Access-Control-Allow-Origin")
	if izin != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, ingin %q", izin, "*")
	}
	metode := w.Header().Get("Access-Control-Allow-Methods")
	if !strings.Contains(metode, "POST") || !strings.Contains(metode, "PUT") {
		t.Errorf("Access-Control-Allow-Methods = %q, harus memuat POST dan PUT", metode)
	}
	if !handlerDipanggil {
		t.Error("handler akhir harus tetap dipanggil untuk GET")
	}
}

// TestCORSPreflight: request OPTIONS harus dijawab 204 dan TIDAK
// meneruskan ke handler (hemat kerja karena preflight tidak butuh data).
func TestCORSPreflight(t *testing.T) {
	handlerDipanggil := false
	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerDipanggil = true
		w.WriteHeader(http.StatusOK)
	})

	app := CORS(dummy)

	req := httptest.NewRequest(http.MethodOptions, "/tasks", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusNoContent)
	}
	if handlerDipanggil {
		t.Error("preflight OPTIONS tidak boleh meneruskan ke handler")
	}
}
