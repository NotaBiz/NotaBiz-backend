package controller

import "net/http"

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("products"))
}