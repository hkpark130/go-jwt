package main

import (
	app "golang/jwt/api"
	"golang/jwt/api/adapter"
	"log"
)

func main() {
	db, err := adapter.Init()
	if err != nil {
		log.Printf("Failed to connect to Database %s ", err)
	}

	r := app.SetupRouter(db)
	r.Run(":3000")
}
