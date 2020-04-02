package main

type route struct {
	auth routerAuth
	main routerMain
	room routerRoom
}

var router = route{}
