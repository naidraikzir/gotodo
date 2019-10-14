package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	uuid "github.com/satori/go.uuid"
)

type Todo struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
}

var (
	db        *gorm.DB
	initDBErr error
	todos     []Todo
)

func initDB() {
	db, initDBErr = gorm.Open("sqlite3", "./todos.db")
	if initDBErr != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Todo{})
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

	allowedMethods := handlers.AllowedMethods([]string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
	})

	// start server
	log.Println("Starting app at :4321")
	log.Fatal(http.ListenAndServe(":4321", handlers.CORS(allowedMethods)(r)))
}

func todosAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo
	db.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

func todosGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo Todo
	id, _ := uuid.FromString(params["id"])
	db.Where("id = ?", id).First(&todo)
	json.NewEncoder(w).Encode(todo)
}

func todosCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	id := uuid.NewV4()
	todo.ID = id
	db.Create(&todo)
	json.NewEncoder(w).Encode(todo)
}

func todosUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var body Todo
	var todo Todo
	id, _ := uuid.FromString(params["id"])
	_ = json.NewDecoder(r.Body).Decode(&body)
	db.Where("id = ?", id).First(&todo)
	todo.Content = body.Content
	todo.Completed = body.Completed
	db.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}

func todosDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := uuid.FromString(params["id"])
	var todo Todo
	db.Where("id = ?", id).Delete(&todo)
	w.Write([]byte(params["id"] + " deleted"))
}
