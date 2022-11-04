package main

import (
	"net/http"

	"github.com/Robster0/fox"
)

func main() {
	r := fox.NewRouter()

	r.GET("/", home)

	r.GET("/profile", auth, profile)

	r.Listen(3000)
}

func auth(c fox.Context) {

}

func home(c fox.Context) {
	c.Status(http.StatusAccepted).Send("<h1>Home Page</h1>")
}

func profile(c fox.Context) {

}
