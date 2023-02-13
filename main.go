package main

//crypto.randomBytes(32).toString('hex')
import (
	"fmt"
	"time"

	"github.com/robin-andreasson/fox"
)

func main() {

	r := fox.Init()

	fox.CORS(fox.CorsOptions{
		Origins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		Methods:     []string{"GET", "POST", "PUT"},
		Headers:     []string{"content-type"},
		Credentials: true,
	})

	fox.Session(fox.SessionOptions{
		Secret:           "tangentbordkatt",
		TimeOut:          1000 * 30,
		ClearProbability: 100,
		Cookie: fox.CookieAttributes{
			HttpOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   1000,
		},
	})

	r.Static("public")

	r.Get("/", index)

	r.Get("/profile/:book;[a-zA-Z0-9]+/:page;[0-9]+", profile)

	validate := r.Group("validate")

	validate.Post("/token", validatetoken)

	auth := r.Group("auth")

	auth.Post("/login", authlogin)

	fmt.Println("Server starting at port 3000")
	r.Listen(3000)
}

func profile(c *fox.Context) {

	c.Text(fox.Status.Ok, fmt.Sprint("book name: ", c.Params["book"], " page is ", c.Params["page"], " desc is ", c.Query["desc"]))
}

func index(c *fox.Context) {

	c.File(fox.Status.Ok, "./html/jwt.html")
}

func validatetoken(c *fox.Context) {

	if c.Session == nil {
		c.JSON(fox.Status.Ok, map[string]string{"error": "no session"})

		return
	}

	c.JSON(fox.Status.Ok, c.Session)
}

func authlogin(c *fox.Context) {

	username := fox.Get[string](c.Body, "user", "person", "username")
	password := fox.Get[int](c.Body, "user", "person", "password")

	if username != "robin" || password != 123 {
		c.JSON(fox.Status.Ok, map[string]any{"error": "wrong username or password"})

		return
	}

	err := c.SetSession(map[string]any{"username": "robin", "password": 123, "iat": time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")})

	if err != nil {
		c.JSON(fox.Status.Ok, map[string]any{"error": err.Error()})
	}

	c.JSON(fox.Status.Ok, map[string]any{"username": "robin", "password": 123})
}
