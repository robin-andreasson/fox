package main

import (
	"fmt"
	"time"

	"github.com/Robster0/fox"
)

func main() {
	r := fox.NewRouter()

	fmt.Println(time.Now().UTC())

	r.GET("/", home)

	r.GET("/profile", auth, profile)

	r.Listen(3000, func(err error) {

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Starting server at port:", r.Port)
	})
}

func auth(c fox.Context) {

}

func home(c fox.Context) {

	c.SetHeader("X-test", "VERY COOL VALUE")

	c.S(fox.Status.Ok, "<h1>Home Page</h1>")
}

func profile(c fox.Context) {

}
