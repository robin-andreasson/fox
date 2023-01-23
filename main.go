package main

import (
	jseon "encoding/json"
	"fmt"
	"log"
	"reflect"

	"os"

	"github.com/robin-andreasson/fox"
)

func main() {
	r := fox.NewRouter()

	r.Static("public")

	r.Get("/", home)

	r.Get("/cookies", cookies)

	r.Get("/profile/:name", auth, profile)

	r.Get("/file", file)

	r.Get("/book/:title;[a-zA-Z]+/:page;[0-9]+", book)

	r.Get("/json", json_get)

	r.Post("/json", json)

	r.Post("/urlencoded", urlencoded)

	r.Post("/image", image)

	fmt.Println("Starting port at", 3000)
	r.Listen(3000)
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

	fmt.Println(JSON_Marshal(c.Json))

	c.JSON(fox.Status.Ok, c.Json.(map[string]any))
}

func urlencoded(c *fox.Context) {

	fmt.Println(c.Form)

	c.Status(fox.Status.Ok)
}

func json_get(c *fox.Context) {
	bytes, _ := os.ReadFile("./data.json")

	mapper := make(map[string]any)
	jseon.Unmarshal(bytes, &mapper)

	c.JSON(fox.Status.Ok, mapper)
}

func image(c *fox.Context) {

	files := fox.Get(c.FormData, "Files", "name-file")

	data := fox.Get(files, "Data").([]byte)

	filename := fox.Get(files, "FileName").(string)

	err := os.WriteFile(filename, data, 0777)

	if err != nil {
		log.Panic(err)
	}

	c.JSON(fox.Status.Ok, c.FormData)
}

func JSON_Marshal(v any) string {
	if !isMap(v) && !isArray(v) {
		return Value(v)
	}

	s := ""

	if isMap(v) {
		s += "{"

		i := 0
		for name, value := range v.(map[string]any) {
			comma := ""

			if i != len(v.(map[string]any)) {
				comma = ","
			}

			var next any

			if isMap(value) {
				next = value.(map[string]any)
			} else {
				next = value
			}

			s +=
				`"` + name + `"` +
					": " +
					JSON_Marshal(next) +
					comma

			i++
		}

		s += "}"
	} else if isArray(v) {
		s += "["

		s += "]"
	}

	return s
}

func isMap(v any) bool {
	if reflect.TypeOf(v).Kind() != reflect.Map {
		return false
	}

	return true
}

func isArray(v any) bool {
	if reflect.TypeOf(v).Kind() != reflect.Array &&
		reflect.TypeOf(v).Kind() != reflect.Slice {
		return false
	}

	return true
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
