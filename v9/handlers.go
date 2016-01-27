package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
    "log"

	"github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
)

//AppContext context
type AppContext struct {
	db *mgo.Database
}


//Index index route
func (c *AppContext) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

//TodoIndex list todos
func (c *AppContext) TodoIndex(w http.ResponseWriter, r *http.Request) {
    repo := TodoRepo{c.db.C("todo")}
    
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
    
    todos, err := repo.RepoAll()
    
    if err == nil {
        if err = json.NewEncoder(w).Encode(todos); err != nil {
		  panic(err)
	   }
    } else {
        log.Println(fmt.Sprintf("error in TodoIndex: %s", err.Error()))
        panic(err)
    }
                   
}

//TodoShow display 1 todo
func (c *AppContext) TodoShow(w http.ResponseWriter, r *http.Request) {
    repo := TodoRepo{c.db.C("todo")}
    
	vars := mux.Vars(r)
	var todoID string
	var err error
	todoID = vars["todoId"]
           
	todo, err := repo.RepoFindTodo(todoID)
	if  err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

/*
TodoCreate with
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/
func (c *AppContext) TodoCreate(w http.ResponseWriter, r *http.Request) {
    repo := TodoRepo{c.db.C("todo")}
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := repo.RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
