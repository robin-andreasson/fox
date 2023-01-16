package main

import (
	"fmt"
	"regexp"

	"os"

	"github.com/robin-andreasson/fox"
)

func main() {
	r := fox.NewRouter()
	hasFile_rex := regexp.MustCompile(`^Content-Disposition: form-data; name="(.+?)"; filename="(.+?)"$`)

	result := hasFile_rex.FindStringSubmatch(`Content-Disposition: form-data; name=""test";"; filename="coolfilename"`)

	for i, re := range result {
		fmt.Println("INDEX: ", i, ": ", re)
	}

	r.Static("public")

	r.Get("/", home)

	r.Get("/profile/:name", auth, profile)

	r.Get("/file", file)

	r.Get("/book/:title;[a-zA-Z]+/:page;[0-9]+", book)

	r.Post("/image", image)

	fmt.Println("Starting port at", 3000)
	r.Listen(3000)
}

func auth(c *fox.Context) {
	fmt.Println("AUTH!")

	c.Cookie("test", "damn thats a good value", fox.CookieAttributes{BASE64: true, ExpiresIn: 1000 * 60 * 60, SameSite: "Lax"})

	c.Next()

}

func image(c *fox.Context) {

	data := c.FormData["Files"].(map[string]interface{})["name-file"].(map[string][]byte)["Data"]

	fmt.Println(data)

	c.String(fox.Status.Created, "test")
}

func home(c *fox.Context) {
	c.String(fox.Status.Ok, "<h1>Home Page</h1>")
}

func profile(c *fox.Context) {

	c.String(fox.Status.Ok, "<h1>"+c.Params["name"]+"'s PROFILE PAGE!</h1>")
}

func file(c *fox.Context) {

	dirname, _ := os.Getwd()

	c.File(fox.Status.Ok, dirname+"/html/index.html")
}

func book(c *fox.Context) {

	c.String(fox.Status.Ok, "<h1>Title: "+c.Params["title"]+" Page: "+c.Params["page"]+"</h1>")
}
