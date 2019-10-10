package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

type Todo struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
}

var (
	db       *sql.DB
	todos    []Todo
	isUuidV4 = uuid.Must(uuid.NewV4())
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Fatal(err)
	}

	stmt, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id TEXT NOT NULL UNIQUE,
			content	TEXT NOT NULL,
			completed INTEGER NOT NULL,
			PRIMARY KEY(id)
		)`,
	)
	stmt.Exec()
}

func main() {
	initDB()

	// init router
	r := mux.NewRouter()
	r.HandleFunc("/todos", todosAll).Methods("GET")
	r.HandleFunc("/todos/{id}", todosGet).Methods("GET")
	r.HandleFunc("/todos", todosCreate).Methods("POST")
	r.HandleFunc("/todos/{id}", todosUpdate).Methods("PUT")
	r.HandleFunc("/todos/{id}", todosDelete).Methods("DELETE")

	// start server
	log.Fatal(http.ListenAndServe(":4321", r))
}

func todosAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo
	rows, _ := db.Query("SELECT * FROM todos")
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Content, &todo.Completed)
		todos = append(todos, todo)
	}
	json.NewEncoder(w).Encode(todos)
}

func todosGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo Todo
	id, _ := uuid.FromString(params["id"])
	row := db.QueryRow("SELECT * FROM todos WHERE id = ?", id)
	err := row.Scan(&todo.ID, &todo.Content, &todo.Completed)
	switch err {
	case sql.ErrNoRows:
		http.Error(w, "Not Found", http.StatusNotFound)
		break
	case nil:
		json.NewEncoder(w).Encode(todo)
		break
	default:
		panic(err)
	}
}

func todosCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	id, _ := uuid.NewV4()
	todo.ID = id
	stmt, _ := db.Prepare("INSERT INTO todos (id, content, completed) VALUES (?, ?, ?)")
	stmt.Exec(todo.ID, todo.Content, todo.Completed)
	stmt.Close()
	json.NewEncoder(w).Encode(todo)
}

func todosUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo Todo
	id, _ := uuid.FromString(params["id"])
	_ = json.NewDecoder(r.Body).Decode(&todo)
	stmt, _ := db.Prepare("UPDATE todos SET content = ?, completed = ? WHERE id = ?")
	stmt.Exec(todo.Content, todo.Completed, id)
	stmt.Close()
	todo.ID = id
	json.NewEncoder(w).Encode(todo)
}

func todosDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := uuid.FromString(params["id"])
	stmt, _ := db.Prepare("DELETE FROM todos WHERE id = ?")
	stmt.Exec(id)
	stmt.Close()
}
