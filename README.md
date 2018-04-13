# Otto
[![GoDoc](https://godoc.org/github.com/JacobSoderblom/otto?status.svg)](http://godoc.org/github.com/JacobSoderblom/otto)
[![Build Status](https://travis-ci.org/JacobSoderblom/otto.svg?branch=master)](https://travis-ci.org/JacobSoderblom/otto)
[![Go Report Card](https://goreportcard.com/badge/github.com/JacobSoderblom/otto)](https://goreportcard.com/report/github.com/JacobSoderblom/otto)

Otto is an easy way to use HTTP in golang. You can either use the Otto App with features like gracefully shutdown on OS signals interrupt and kill, or just use the Otto Router and handle the HTTP Server as you like!

Otto uses [Gorilla/mux](https://github.com/gorilla/mux) to define routes.

## Feature overview

- Group APIs
- Middleware
- Functions that makes it easy to send HTTP responses
- Centralized HTTP error handling
- Custom error handlers to specific HTTP status codes
- Possibility to only use the router part

## Examples

Example of how to use the Otto App

```Go
package main

import (
  "log"
  "github.com/JacobSoderblom/otto"
)

func main() {
  opts := otto.NewOptions()
  opts.Addr = ":3000"
  
  app := otto.New(opts)
  
  app.GET("/", func(ctx otto.Context) error {
    return ctx.String("Hello world!")
  })
  
  log.Fatal(app.Serve())
}
```

Example of how to use the Router without the Otto App.

The Router implements the `http.Handler` so it is easy to use without the `App`!


```Go
package main

import (
  "log"
  "github.com/JacobSoderblom/otto"
)

func main() {
  r := otto.New(opts)
  
  r.GET("/", func(ctx otto.Context) error {
    return ctx.String("Hello world!")
  })
  
  log.Fatal(http.ServeAndListen(":3000", r)
}
```

Example of using middlewares


```Go
package main

import (
  "log"
  "github.com/JacobSoderblom/otto"
)

func main() {
  r := otto.New(opts)
  
  r.Use(func(next otto.HandlerFunc) otto.HandlerFunc {
    return func(ctx otto.Context) error {
      ctx.Set("some key", "I'm a middleware!")
      return next(ctx)
    }
  })
  
  r.GET("/", func(ctx otto.Context) error {
    return ctx.String("Hello world!")
  })
  
  log.Fatal(http.ServeAndListen(":3000", r)
}
```

Example of how to associate error handler to HTTP status code


```Go
package main

import (
  "log"
  "github.com/JacobSoderblom/otto"
)

func main() {
  r := otto.New(opts)
  
  errorHandlers := map[int]otto.ErrorHandler{
    400: func(code int, err error, ctx otto.Context) error {
      // do some awesome error handling!
      return ctx.Error(code, err)
    },
    401: func(code int, err error, ctx otto.Context) error {
      // do some awesome error handling!
      return ctx.Error(code, err)
    },
  }
  
  r.SetErrorHandlers(errorHandlers)
  
  r.GET("/", func(ctx otto.Context) error {
    return ctx.String("Hello world!")
  })
  
  log.Fatal(http.ServeAndListen(":3000", r)
}
```
