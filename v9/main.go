package main

import (
	"log"
	"net/http"
    "gopkg.in/mgo.v2"
)

func main() {

	session, err := mgo.Dial("10.10.12.128")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := AppContext{session.DB("test")}

	router := NewRouter(&appC)

	log.Fatal(http.ListenAndServe(":8080", router))
}
