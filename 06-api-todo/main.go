package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Todo = model data kita (sama seperti latihan di langkah 4)
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// "database" in-memory: slice of struct
var todos = []Todo{
	{ID: 1, Title: "Belajar variabel", Done: true},
	{ID: 2, Title: "Belajar fungsi", Done: true},
	{ID: 3, Title: "Belajar loop", Done: true},
}

var nextID = 4 // counter untuk ID baru

// ---- Handler: GET /todos ----
func listTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

// ---- Handler: POST /todos ----
func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	// Baca body JSON dari request
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "format JSON tidak valid", http.StatusBadRequest)
		return
	}
	// Validasi: title tidak boleh kosong
	if strings.TrimSpace(newTodo.Title) == "" {
		http.Error(w, "title tidak boleh kosong", http.StatusBadRequest)
		return
	}
	newTodo.ID = nextID
	nextID++
	todos = append(todos, newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

// ---- Handler: DELETE /todos/{id} ----
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// Ambil {id} dari URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID harus angka", http.StatusBadRequest)
		return
	}

	// Cari todo dan hapus dari slice
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "todo tidak ditemukan", http.StatusNotFound)
}

// ---- Handler: PATCH /todos/{id} (tandai selesai) ----
func toggleTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID harus angka", http.StatusBadRequest)
		return
	}

	// Cari todo, ubah Done jadi true
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Done = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "todo tidak ditemukan", http.StatusNotFound)
}

func main() {
	// Daftarkan route. Pola "METHOD /path" berlaku di Go 1.22+
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", listTodos)
	mux.HandleFunc("POST /todos", createTodo)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodo)
	mux.HandleFunc("PATCH /todos/{id}", toggleTodo)

	log.Println("Server jalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
