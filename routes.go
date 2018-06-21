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
		"/api/user/{id}",
		getUserHandler,
	},
	route{
		"User_POST",
		"POST",
		"/api/user",
		postUserHandler,
	},
	route{
		"User_PUT",
		"PUT",
		"/api/user/{id}",
		putUserHander,
	},
	route{
		"User_DELETE",
		"DELETE",
		"/api/user/{id}",
		deleteUserHander,
	},
	route{
		"Users_GET",
		"GET",
		"/api/users",
		getUsersHandler,
	},
}
