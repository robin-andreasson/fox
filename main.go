package main

//crypto.randomBytes(32).toString('hex')
import (
	"fmt"
	"log"
	"os"

	"github.com/robin-andreasson/fox"
)

func main() {

	r := fox.Init()

	fox.CORS(fox.CorsOptions{
		Origins:     []string{"http://127.0.0.1:5500", "*"},
		Methods:     []string{"POST", "*"},
		Headers:     []string{"content-type"},
		Credentials: true,
	})

	r.Static("public")

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
	c.Cookie("token", "this is a epic token value", fox.CookieAttributes{BASE64: true, MaxAge: 60 * 60 * 24})

	c.JSON(fox.Status.Ok, c.ResHeaders())
}

func auth(c *fox.Context) {
	fmt.Println("AUTH!")

	c.Next()
}

func json(c *fox.Context) {

	c.Cookie("token", "This is an insane token value", fox.CookieAttributes{
		BASE64:   true,
		HttpOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
		MaxAge:   60 * 60 * 24,
	})

	c.JSON(fox.Status.Ok, c.Body)
}

func urlencoded(c *fox.Context) {

	fmt.Println(c.Body)

	firstname := fox.Get[string](c.Body, "firstname")
	lastname := fox.Get[string](c.Body, "lastname")

	fmt.Println(firstname)
	fmt.Println(lastname)

	c.JSON(fox.Status.Ok, c.Body)
}

func json_get(c *fox.Context) {

	c.JSON(fox.Status.Ok, c.Headers)
}

func image(c *fox.Context) {

	//fmt.Println(c.Body)

	files := fox.Get[map[string]any](c.Body, "Files", "image")

	data := fox.Get[[]byte](files, "Data")

	filename := fox.Get[string](files, "Filename")

	err := os.WriteFile(filename, data, 0777)

	if err != nil {
		log.Panic(err)
	}

	c.File(fox.Status.Ok, filename)
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

	c.File(fox.Status.Ok, "./html/index.html")
}

func book(c *fox.Context) {

	c.Text(fox.Status.Ok, "<h1>Title: "+c.Params["title"]+" Page: "+c.Params["page"]+"</h1>")
}
