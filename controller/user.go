package controller

import "net/http"

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get users"))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register users"))
}