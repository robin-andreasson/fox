# Fox

Web Application Framework built in Go.

Session feature uses sqlite as session storage, you'll need an [gcc compiler](https://jmeubank.github.io/tdm-gcc) to run your program.

# Layout

```go
//Initialize root router
r := fox.Init()

/*
    ...routes and middlewares
*/

//Start server at port 3000
r.Listen(3000)
```

# Features


#### Routing
Map functions to url paths
```go
r.Get("/", handler)

r.Post("/", handler)

r.Put("/", handler)

r.Delete("/", handler)
```
#### Groups
Add path prefixes to remove redundancy
```go
auth := r.Group("/auth")

// -> "/auth/signin"
auth.Get("/signin", handler)
```

#### Path variables

Add variables in path, colon as delimiter
```go
//Variable named id
r.Get("/task/:id", getTask)

func getTask(c *fox.Context) error {
    //get id
    id := c.Params("id")
}
```
Regex pattern for validation, semicolon as delimiter
```go
/*
Trigger handler if id variable matches the numbers only pattern

Fox automatically wraps your regex statement between ^ and $
e.g. \d+ becomes ^\d+$
*/
r.Get("/user/:id;\d+")
```

#### Body parsing and fox.Get
Fox parses your http body and maps them to an interface typed struct field in fox.Context called Body.


Multipart/form-data example:
```go
/*
c.Body: 
{
    "Files": {
        "file-name": {
            "Data": []byte
            "Filename": string
            "Content-Type": string
        }
    }
}
*/
r.Post("/image", handler)

func handler(c *fox.Context) error {

    //retrieve nested interface data with fox.Get
    image := fox.Get[map[string]any](c.Body, "Files", "file-name")

    data := fox.Get[[]byte](image, "Data")
    filename := fox.Get[string](image, "Filename")

    if err := os.WriteFile(filename, data, 777); err != nil {
        return c.Status(fox.Status.InternalServerError)
    }

    return c.Status(fox.Status.Created)
}
```

#### Session

Session middleware, a way to store information to be used across multiple pages.

```go

/*
Initialize session options
*/
fox.Session(fox.SessionOptions{
    Secret:           "YOUR_SECRET",
    TimeOut:          1000 * 60 * 60 * 24 * 30,
    ClearProbability: 2.5,
    Path: "./storage/session.sql",

    Cookie: fox.CookieAttributes{
    	HttpOnly:  true,
    	Secure:    true,
    	SameSite:  "Lax",
    	Path:      "/",
    	ExpiresIn: 1000 * 60 * 60 * 24 * 30,
    },
})


/*
Set session
*/
r.Post("/auth/signin", signin)

func signin(c *fox.Context) error {
    //store information and create FOXSESSID cookie
    c.SetSession(map[string]string{/*...data*/})
 
    c.Status(fox.Status.Ok)
}


/*
Get session
*/
r.Post("/session", session)

func session(c *fox.Context) error {
    c.Session// -> map[string]string{/*...data*/}
}
```

#### Refresh

allows the client to continuously receive new jwt access token through the *_X-Fox-Access-Token_* header.


```go
//initialize refresh options
fox.Refresh(fox.RefreshOptions{
    //Found in X-Fox-Access-Token header
    AccessToken: fox.TokenOptions{
    	Secret: "YOUR_ACCESS_TOKEN_SECRET",
    	Exp:    1000 * 60 * 5,
    },
    //Stored in FOXREFRESH cookie
    RefreshToken: fox.TokenOptions{
    	Secret: "YOUR_REFRESH_TOKEN_SECRET",
    	Exp:    1000 * 60 * 60 * 24,
    },
    Cookie: fox.CookieAttributes{
    	HttpOnly:  true,
    	Secure:    true,
    	SameSite:  "Lax",
    	Path:      "/",
    	ExpiresIn: 1000 * 60 * 60 * 24,
    },
    RefreshFunction: handleRefresh,
})

/*
used to fetch the new access token data after it has expired.

refreshdata parameter is the data stored in refresh token.
*/
func handleRefresh(refreshdata any) (any, error) {

    id := fox.Get[string](refreshdata, "id")

    result := /* fetch */

    return result, nil
}


r.Post("/auth/signin", signin)

func signin(c *fox.Context) error {

    //add non vulnerable data for fetching new access token data
    refreshData := map[string]string{/*...data */}

    accessData := map[string]string{/*...data */}
    
    c.SetRefresh(accessData, refreshData)
 
    c.Status(fox.Status.Ok)
}

/*
Get refresh
*/
r.Post("/refresh", session)

func session(c *fox.Context) error {
    c.Refresh// -> map[string]string{/*...data*/}
}
```

#### Middleware

Add functions to be called before the main handler

```go

r.Post("/", middleware, handler)

func middleware(c *fox.Context) error {
    //continue to the next middleware/handler in the stack
    return c.Next()    
}
```

#### Statically serve files

Specify folder used to serve static files when the client endpoint requests them
```go
/*
    wd: /src/internal
*/

//Serve folder named public in /src/internal 
r.Static("public")

//Serve folder named static in /src
r.Static("static", "../")

```