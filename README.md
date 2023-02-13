# Fox

Go package for handling http requests  


# Features

* Basic routing, e.g. Get, "/home"
* url params, e.g "/profile/:name;[a-zA-Z]+"
* Body parsing (json, formdata, urlencoded (with nested keys))
* fox.Get to obtain data from nested map interface
* Statically serve files
* Create route groups
* Session and Refresh middlewares 


```go
//Initialize a new router
r := fox.Init()

//Basic Get handler
r.Get("/", handler)

//Matches with any url path that begins with /
r.Get("/*", handler)
//Matches with any url path that begins with /code-
//e.g. /code-418 or /code-200/example
r.Get("/code-*", handler)

//Middleware stack
r.Get("/auth", auth_mw, handler)

func auth_mw(c *fox.Context) error {
    //Continue to the next handler inside the stack
    return c.Next()
}

//Params, delimiter: " : "
r.Get("/post/:id", func(c *fox.Context) error {

    c.Params["id"]
})

//Regex pattern for params, Delimiter: " ; "
r.Get("/book/:title;[a-zA-Z]+/:page;[0-9]+", func (c *fox.Context) error {
    //a-z or A-Z
    c.Params["title"]

    //0-9
    c.Params["page"]
})

//Create groups, gets the /api prefix
api := r.Group("api")

/*    /api/json   */
api.Get("/json", handler)



//fox.Get example with a formdata request:
/*
    c.Body:
    {
        "Files": {
            "key-name": {
                "Data": []byte
                "Filename": string
                "Content-Type": string
            }
        }
    }
*/
r.Post("/image", image_handler)

func image_handler(c *fox.Context) error {

    //fox.Get gives you the ability to get values inside a nested map interface easily
    name := fox.Get[string](c.Body, "Files", "key-name")

	data := fox.Get[[]byte](name, "Data")
	filename := fox.Get[string](name, "Filename")

	if err := os.WriteFile(filename, data, 777); err != nil {
        return err
    }

    //Send back status code 201
	return c.Status(fox.Status.Created)
}


//Start server at specified port
r.Listen(3000)
```