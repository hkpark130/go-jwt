package main

import (
	"golang/jwt/app"
)

func main() {
	r := app.SetupRouter()
	r.Run(":3000")
}
