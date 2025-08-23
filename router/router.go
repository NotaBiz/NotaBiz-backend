package router

import (
	"NotaBiz-backend/controller"
	"net/http"
)

func InitRouter() *http.ServeMux {
	mux := http.NewServeMux()
	apiGroup := map[string](map[string]http.HandlerFunc){
		"/users":    usersHandler,
		"/products": productsHandler,
	}

	for path, api := range apiGroup {
		for method, handler := range api {
			mux.HandleFunc("/api/v1"+path+method, handler)
		}
	}
	return mux
}

var usersHandler = map[string]http.HandlerFunc{
	"/":         controller.GetUsers,
	"/register": controller.RegisterUser,
}

var productsHandler = map[string]http.HandlerFunc{
	"/": controller.GetProducts,
}
