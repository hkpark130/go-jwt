package main

import (
	"golang/jwt/app"
	"log"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	r := app.SetupRouter(path) // 「go run」の時、パスが（go-jwt）になる
	r.Run(":3000")
}
