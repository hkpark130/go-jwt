package main

import (
	"golang/jwt/app"
)

func main() {
	path := "./"
	r := app.SetupRouter(path) // 「go run」の時、パスが（go-jwt）になる

	r.Run(":3000")
}
