package main

import (
	"fmt"
	"strings"

	"github.com/gocraft/web"
)

type Context struct {
	HelloCount int
}

func (c *Context) SetHelloCount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.HelloCount = 3
	next(rw, req)
}

func (c *Context) SayHello(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, strings.Repeat("Hello ", c.HelloCount), "World!")
}

func main() {

	// router := web.New(Context{}). // Create your router
	// 				Middleware(web.LoggerMiddleware).    // Use some included middleware
	// 				Middleware(web.ShowErrorsMiddleware) // ...

	// router.Get("/", (*Context).SayHello)
	// router.Get("/", (*Context).SayHello)
	// http.ListenAndServe("localhost:3000", router) // Start the server!
}
