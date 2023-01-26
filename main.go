package main

import (
	jseon "encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/robin-andreasson/fox"
)

const (
	rex_s_str = "&|="
	rex_n_str = `\[(.+?)\]`
)

func Urlencoded(s string) map[string]any {
	body := make(map[string]any)

	rex_s := regexp.MustCompile(rex_s_str)
	rex_n := regexp.MustCompile(rex_n_str)

	seg := rex_s.Split(s, -1)

	for i := 0; i < len(seg); i += 2 {

		name := seg[i]
		value := seg[i+1]

		nestedkeys := rex_n.FindAllStringSubmatch(name, -1)

		//if there are nested keys
		if len(nestedkeys) != 0 {

			n, _, _ := strings.Cut(name, "[")

			if body[n] == nil {
				body[n] = make(map[string]any)
			}

			body[n] = getNestedKeys(nestedkeys, value, body[n])

			continue
		}

		body[name] = value
	}

	return body
}

func getNestedKeys(names [][]string, value string, body any) any {
	if len(names) == 0 {
		return value
	}

	name := names[0][1]
	names = names[1:]

	next := body.(map[string]any)

	if next[name] == nil || reflect.TypeOf(next[name]).Kind() != reflect.Map {
		next[name] = make(map[string]any)
	}

	nested_value := getNestedKeys(names, value, next[name])

	next[name] = nested_value

	return body
}

func main() {
	r := fox.NewRouter()

	s := Urlencoded("person[firstname]=robin&person[lastname]=andreasson")

	fmt.Println(s)

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

	fmt.Println(c.Body)

	c.JSON(fox.Status.Ok, c.Body.(map[string]any))
}

func urlencoded(c *fox.Context) {

	c.Status(fox.Status.Ok)
}

func json_get(c *fox.Context) {
	bytes, _ := os.ReadFile("./data.json")

	mapperd := make(map[string]any)
	jseon.Unmarshal(bytes, &mapperd)

	c.JSON(fox.Status.Ok, mapperd)
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
