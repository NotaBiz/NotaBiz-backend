// Package main implements the entry point for the NotaBiz-backend application.
// It sets the local time zone to "Asia/Jakarta" and starts the server.
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