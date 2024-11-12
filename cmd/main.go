package main

import (
	"test-server/internal/app"
	"test-server/internal/server"
)

const connStr = "postgres://postgres:1234@0.0.0.0:5432/shop_dp?sslmode=disable"

func main() {

	serviceApp := app.NewApp(connStr)
	srv := server.New(serviceApp)
	srv.Run()

}
