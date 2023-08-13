# Glam

A functional, composible http router for go. Fast, zero alloc*, and all around badass and glamorous

## Features

Glam is a minimal go http router and not bloated with tons of extra features. Check out the minimal set of features below

### Strict Route Handlers

Strict routes handlers handle only an exact match of the url path. They have the highest precidence of any route handler. See the precidence section below for more details

```go

    router := glam.NewRouter()
    //This route matches the exact path '/strict'
    router.Get("/strict", func(w http.ResponseWriter, r *http.Request){
        ...
    })
```

### Static Route Handlers

Static route handlers handle

## Route Handler Precidence

1) Strict 
2) Static
3) Regex
4) Param
5) Not Found