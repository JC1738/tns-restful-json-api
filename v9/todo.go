package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Todo struct
type Todo struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Completed bool          `json:"completed" bson:"completed"`
	Due       time.Time     `json:"due" bson:"due"`
}

//TodoCollection collection
type TodoCollection struct {
	Data []Todo `json:"data"`
}

//TodoResource resource
type TodoResource struct {
	Data Todo `json:"data"`
}


//Todos array
type Todos []Todo
