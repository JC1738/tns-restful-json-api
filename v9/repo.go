package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//TodoRepo mongo db
type TodoRepo struct {
	coll *mgo.Collection
}

//Init Give us some seed data
func (r *TodoRepo) Init() {
	r.RepoCreateTodo(Todo{Name: "Write presentation"})
	r.RepoCreateTodo(Todo{Name: "Host meetup"})
}

//RepoFindTodo find
func (r *TodoRepo) RepoFindTodo(id string) (Todo, error) {
	result := TodoResource{}

	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	if err != nil {
		log.Println(fmt.Sprintf("error in RepoFindTodo: %s", err.Error()))
		return result.Data, err
	}

	log.Println(fmt.Sprintf("RepoFindTodo: found %s with data %s", id, result.Data))

	return result.Data, nil
}

//RepoAll return all
func (r *TodoRepo) RepoAll() (TodoCollection, error) {
	result := TodoCollection{[]Todo{}}
	err := r.coll.Find(nil).All(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

//RepoCreateTodo this is bad, I don't think it passes race condtions
func (r *TodoRepo) RepoCreateTodo(t Todo) Todo {

	id := bson.NewObjectId()
	_, err := r.coll.UpsertId(id, t)
	if err != nil {
		//handle err
	}

	t.ID = id

	return t
}

//RepoDestroyTodo delete
func (r *TodoRepo) RepoDestroyTodo(id string) error {
	err := r.coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	return nil
}
