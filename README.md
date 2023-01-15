# Fox

Go package for handling http requests  


# Features

* Basic routing, e.g. Get, "/home" etc
* Querystrings, url params (with custom regex checking)



```go
//Initialize a new router
r := fox.NewRouter()

//Basic Get handler
r.Get("/", handler)

//Middleware stack
r.Get("/auth", auth_mw, handler)

//Params, delimiter: " : "
r.Get("/post/:id", func(c *fox.Context) {

    c.Params["id"]
})

//Regex pattern for the params, Delimiter: " ; "
r.Get("/book/:title;[a-zA-Z]+/:page;[0-9]+", handler)


//Start server at specified port
r.Listen(3000)
```