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

	fox.Refresh(fox.RefreshOptions{
		Secret: "tangentbordkatt",
		RefreshFunction: func(refreshobj any) (any, error) {

			return map[string]any{"username": "robin", "password": 123, "iat": time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")}, nil
		},
		AccessToken: fox.TokenOptions{
			Exp: 1000 * 30,
		},
		RefreshToken: fox.TokenOptions{
			Exp: 1000 * 60 * 60 * 24,
		},
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

	validate := r.Group("validate")

	validate.Post("/token", validatetoken)

	auth := r.Group("auth")

	auth.Post("/login", authlogin)

	fmt.Println("Server starting at port 3000")
	r.Listen(3000)
}

func index(c *fox.Context) error {

	return c.File(fox.Status.Ok, "./html/jwt.html")
}

func validatetoken(c *fox.Context) error {

	fmt.Println(c.Refresh)

	return c.JSON(fox.Status.Ok, c.Refresh)
}

func authlogin(c *fox.Context) error {

	username := fox.Get[string](c.Body, "user", "person", "username")
	password := fox.Get[int](c.Body, "user", "person", "password")

	if username != "robin" || password != 123 {
		return c.JSON(fox.Status.Ok, map[string]any{"error": "wrong username or password"})
	}

	accesstoken, err := c.SetRefresh(map[string]any{"username": "robin", "password": 123, "iat": time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")}, map[string]int{"user-id": 1})

	if err != nil {
		return c.JSON(fox.Status.Ok, map[string]any{"error": err.Error()})
	}

	return c.JSON(fox.Status.Ok, map[string]any{"username": "robin", "password": 123, "accesstoken": accesstoken})
}
