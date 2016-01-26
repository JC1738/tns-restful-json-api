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
	h http.Handler
}

func (wrap *WrapHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggedResponse{ResponseWriter: w, status: 200}

	elapsed := callServer(wrap, lw, r)

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
func callServer(w *WrapHTTPHandler, lw *loggedResponse, r *http.Request) (elapsed time.Duration) {
	start := time.Now()

	defer func() {
		elapsed = time.Since(start)
	}()

	w.h.ServeHTTP(lw, r)

	return elapsed
}

//NewRouter setup and return the router for program
func NewRouter(appC *AppContext) *gormux.Router {

	router := gormux.NewRouter().StrictSlash(true)

    rc := new(RoutesCollection)
    rc.BuildRoute(appC)

	commonHandlers := alice.New(g.CompressHandler)

	for _, route := range rc.Routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(&WrapHTTPHandler{h: commonHandlers.Then(handler)})
	}

	return router
}
