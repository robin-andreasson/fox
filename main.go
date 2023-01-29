package main

import (
	jseon "encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/robin-andreasson/fox"
)

func main() {

	r := fox.Root()

	r.Static("public")

	maskedData := 255 & 0xFF

	fmt.Println(maskedData)

	fmt.Println(0x52)

	r.Get("/", home)

	r.Get("/cookies", cookies)

	r.Get("/profile/:name", auth, profile)

	r.Get("/file", file)

	r.Get("/book/:title;[a-zA-Z]+/:page;[0-9]+", book)

	api := r.Group("api")

	api.Get("/json", json_get)
	api.Post("/json", json)

	form := r.Group("form")

	form.Post("/urlencoded", urlencoded)
	form.Post("/image", image)

	fmt.Println("Starting port at", 3000)

	fox.Listen(r, 3000)
}

func cookies(c *fox.Context) {
	fmt.Println(c.Headers["Cookie"])

	fmt.Println(c.Cookies)

	c.Status(fox.Status.ImaTeapot)
}

func auth(c *fox.Context) {
	fmt.Println("AUTH!")

	c.Next()
}

func json(c *fox.Context) {

	//fmt.Println(c.Body)
	//
	//firstname := fox.Get(c.Body, "person", "firstname")
	//lastname := fox.Get(c.Body, "person", "lastname")
	//
	//arr := fox.Get(c.Body, "test", "arr")
	//
	//fmt.Println(firstname)
	//fmt.Println(lastname)
	//fmt.Println(arr)

	c.JSON(fox.Status.Ok, c.Body)
}

func urlencoded(c *fox.Context) {

	fmt.Println(c.Body)

	firstname := fox.Get(c.Body, "firstname")
	lastname := fox.Get(c.Body, "lastname")

	fmt.Println(firstname)
	fmt.Println(lastname)

	c.JSON(fox.Status.Ok, c.Body)
}

func json_get(c *fox.Context) {
	bytes, _ := os.ReadFile("./data.json")

	mapperd := make(map[string]any)
	jseon.Unmarshal(bytes, &mapperd)

	c.JSON(fox.Status.Ok, mapperd)
}

func image(c *fox.Context) {

	files := fox.Get(c.Body, "Files", "post-image")

	data := fox.Get(files, "Data").([]byte)

	filename := fox.Get(files, "FileName").(string)

	err := os.WriteFile(filename, data, 0777)

	if err != nil {
		log.Panic(err)
	}

	c.File(fox.Status.Ok, "./"+filename)
}

func home(c *fox.Context) {
	c.Text(fox.Status.Ok, "<h1>Home Page</h1>")
}

func profile(c *fox.Context) {

	c.Text(fox.Status.Ok, "<h1>"+c.Params["name"]+"'s PROFILE PAGE!</h1>")
}

func file(c *fox.Context) {

	c.Cookie("test", "damn thats a good value", fox.CookieAttributes{BASE64: true, ExpiresIn: 1000 * 60 * 60, SameSite: "Lax"})
	c.Cookie("name", "BAD ; VALUE", fox.CookieAttributes{ExpiresIn: 1000 * 60 * 60, SameSite: "Lax"})
	c.Cookie("test2", "DAMN, GOOD VALUE", fox.CookieAttributes{ExpiresIn: 1000 * 60 * 60, SameSite: "Strict"})

	dirname, _ := os.Getwd()

	c.File(fox.Status.Ok, dirname+"/html/index.html")
}

func book(c *fox.Context) {

	c.Text(fox.Status.Ok, "<h1>Title: "+c.Params["title"]+" Page: "+c.Params["page"]+"</h1>")
}
