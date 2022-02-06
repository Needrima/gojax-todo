package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type todo struct {
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

var todos = map[int]todo{} // in memory database

func init() {
	rand.Seed(time.Now().UnixNano())
	todos = map[int]todo{
		rand.Intn(30): {Description: "Go watch football", Status: true},
		rand.Intn(30): {Description: "Go play video games"},
	}
}

func updateTodo(id int, text string, td map[int]todo) {
	switch text {
	case "Done":
		todo := td[id]
		todo.Status = true
		td[id] = todo
	case "Undo":
		todo := td[id]
		todo.Status = false
		td[id] = todo
	case "Cancle":
		delete(todos, id)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/index.html")
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(todos)
		case http.MethodPost:
			var newTodo todo
			json.NewDecoder(r.Body).Decode(&newTodo)
			todos[rand.Intn(30)] = newTodo
		case http.MethodDelete:
			todos = map[int]todo{}
		case http.MethodPut:
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Readall:", err)
				return
			}

			parts := strings.Split(string(bs), " ")

			id, text := parts[0], parts[1]

			idInt, err := strconv.Atoi(id)
			if err != nil {
				log.Println("Number not an integer")
				return
			}
			updateTodo(idInt, text, todos)
		default:
			http.Error(w, "Status method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
