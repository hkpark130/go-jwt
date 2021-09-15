package main

import (
	"golang/jwt/app"
	"os"
)

func main() {
	path, _ := os.Getwd()
	r := app.SetupRouter(path) // 「go run」の時、パスが（go-jwt）になる
	r.Run(":3000")
}
