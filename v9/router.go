package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	g "github.com/gorilla/handlers"
	gormux "github.com/gorilla/mux"
	"github.com/justinas/alice"
)

//WrapHTTPHandler to capture status code, wrapping handler
type WrapHTTPHandler struct {
	m http.Handler
}

func (h *WrapHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggedResponse{ResponseWriter: w, status: 200}

	elapsed := callServer(h, lw, r)

	if lw.status != http.StatusOK {
		log.Println(fmt.Sprintf("ServeHTTP [%s] %s - %d\n", r.RemoteAddr, r.URL, lw.status))
	} else {
		log.Println(fmt.Sprintf("ServeHTTP [%s] %s - %d time: %v\n", r.RemoteAddr, r.URL, lw.status, elapsed))
	}
}

type loggedResponse struct {
	http.ResponseWriter
	status int
}

func (l *loggedResponse) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

//created to measure duration of call to ServeHTTP
func callServer(h *WrapHTTPHandler, lw *loggedResponse, r *http.Request) (elapsed time.Duration) {
	start := time.Now()

	defer func() {
		elapsed = time.Since(start)
	}()

	h.m.ServeHTTP(lw, r)

	return elapsed
}

//NewRouter setup and return the router for program
func NewRouter() *gormux.Router {

	router := gormux.NewRouter().StrictSlash(true)

	commonHandlers := alice.New(g.CompressHandler)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(&WrapHTTPHandler{m: commonHandlers.Then(handler)})
	}

	return router
}
