package main

import "net/http"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routeArray []route

var routes = routeArray{
	route{
		"User_GET",
		"GET",
		"/user/{id}",
		getUserHandler,
	},
	route{
		"User_POST",
		"POST",
		"/user",
		postUserHandler,
	},
	route{
		"User_PUT",
		"PUT",
		"/user/{id}",
		putUserHander,
	},
	route{
		"User_DELETE",
		"DELETE",
		"/user/{id}",
		deleteUserHander,
	},
	route{
		"Users_GET",
		"GET",
		"/users",
		getUsersHandler,
	},
}
