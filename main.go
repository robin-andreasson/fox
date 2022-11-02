package main

import (
	"fmt"

	"github.com/Robster0/fox"
)

func main() {
	fmt.Println(fox.Test)

	r := fox.NewRouter()

	r.GET("/", home)

	r.GET("/profile", auth, profile)

	fmt.Println(r)
}

func auth(c fox.Context) {

}

func home(c fox.Context) {

}

func profile(c fox.Context) {

}
