package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"User_GET",
		"GET",
		"/user/{id}",
		User_GET,
	},
	Route{
		"User_POST",
		"POST",
		"/user",
		User_POST,
	},
	Route{
		"User_PUT",
		"PUT",
		"/user/{id}",
		User_PUT,
	},
	Route{
		"User_DELETE",
		"DELETE",
		"/user/{id}",
		User_DELETE,
	},
	Route{
		"Users_GET",
		"GET",
		"/users",
		Users_GET,
	},
}
