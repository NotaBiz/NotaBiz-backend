package main

import (
	"NotaBiz-backend/app"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	app.RunServer()
}