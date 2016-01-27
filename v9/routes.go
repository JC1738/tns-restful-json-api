package main

import "net/http"

//Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//RoutesCollection collection
type RoutesCollection struct {
    Routes []Route
}

//BuildRoute set the route for the application
func (rc *RoutesCollection) BuildRoute (appc *AppContext)  {

    var routes []Route 

    routes = []Route{
        Route{
            "Index",
            "GET",
            "/",
            appc.Index,  //ok
        },
        Route{
            "TodoIndex",
            "GET",
            "/todos",
            appc.TodoIndex,  //ok
        },
        Route{
            "TodoCreate",
            "POST",
            "/todos",
            appc.TodoCreate,  //ok
        },
        Route{
            "TodoShow",
            "GET",
            "/todos/{todoId}",
            appc.TodoShow,  //ok
        },
    }
    
    rc.Routes = routes    
}


