package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test handler listTodos: pastikan status 200 dan response berupa list JSON
func TestListTodos(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	listTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusOK)
	}

	var hasil []Todo
	if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
		t.Fatalf("response bukan JSON yang valid: %v", err)
	}
	if len(hasil) == 0 {
		t.Error("list kosong, harusnya berisi data awal")
	}
}

// Test handler createTodo: status 201 dan data tersimpan
func TestCreateTodo(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title":"Test tambah"}`))
	w := httptest.NewRecorder()

	createTodo(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusCreated)
	}

	var dibuat Todo
	if err := json.Unmarshal(w.Body.Bytes(), &dibuat); err != nil {
		t.Fatalf("response bukan JSON: %v", err)
	}
	if dibuat.Title != "Test tambah" {
		t.Errorf("title = %q, ingin %q", dibuat.Title, "Test tambah")
	}
	if dibuat.ID == 0 {
		t.Error("ID tidak boleh 0")
	}
}

// Test handler createTodo dengan title kosong: harus 400
func TestCreateTodoTitleKosong(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title":""}`))
	w := httptest.NewRecorder()

	createTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusBadRequest)
	}
}

// Test handler deleteTodo: status 204 saat berhasil
func TestDeleteTodo(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	deleteTodo(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusNoContent)
	}
}

// Test handler deleteTodo dengan ID yang tidak ada: harus 404
func TestDeleteTodoTidakAda(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	deleteTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusNotFound)
	}
}

// Test handler toggleTodo: Done berubah jadi true
func TestToggleTodo(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/todos/2", nil)
	req.SetPathValue("id", "2")
	w := httptest.NewRecorder()

	toggleTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("kode = %d, ingin %d", w.Code, http.StatusOK)
	}

	var hasil Todo
	if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
		t.Fatalf("response bukan JSON: %v", err)
	}
	if !hasil.Done {
		t.Error("Done harusnya true setelah PATCH")
	}
}
