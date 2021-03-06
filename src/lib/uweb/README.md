#uweb Web Framework

uweb is a micro web kit written in Golang. 
It borrows some ideas from Koa.js, Gin, Playframework, beego, etc.

## example
```
//
// src/app/main.go
// 
package main

import (
	"github.com/ot24net/uweb"
	
	_ "ctrls/account"
    _ "models/account"
)

func main() {
	// app
	app := uweb.NewApp()
	
	// Ignore some path
	app.Use(uweb.MdIgnore([]string{"/hacheck"}))
	
	// Response favicon 
	app.Use(uweb.MdFavicon("../../pub/img/favicon.ico"))
	
	// Serve static files, "/pub" is path prefix, and "../../pub" is file directory
	app.Use(uweb.MdStatic("/pub", "../../pub")) // before compress
	
	// Compress use gzip, not work with MdSatic
	app.Use(uweb.MdCompress())
	
	// log
	app.Use(uweb.MdLogger(uweb.LOG_LEVEL_2))
	
	// Cache use memcache
	app.Use(uweb.MdCache("memcache", "127.0.0.1:11000", "cache_prefix"))
	
	// Session depends on cache
	app.Use(uweb.MdSession(3600*12))
	
	// Flash depends on session
	app.Use(uweb.MdFlash())
	
	// Csrf depends on session, and get the Csrf token from session with key: uweb.CSRF_TOKEN_KEY
	app.Use(uweb.MdCsrf())
	
	// Html render
	app.Use(uweb.MdRender("../../pub/html", ".html"))
	
	// Cors
	app.Use(uweb.MdCors(uweb.DefaultCors))
	
	// I18n, depends on session if detect is true
	app.Use(uweb.MdI18n("../../pub/locale", "zh_cn", false))

	// errors page
	app.Use(uweb.MdErrPage(uweb.Map{
		"404_leave_url": "http://baidu.com",
	}))
	
	// pjax
	app.Use(uweb.MdPjax())
	
	// if you want more method, change route.go
	app.Use(uweb.MdRouter())
	
	// listen address
	app.Listen(":9099")
}

//
// src/ctrls/account/login.go 
//
package account

import (
	   "github.com/ot24net/uweb"
	   "models/account"
)

func init() {
	 // simple get
	 uweb.Get("/account/login", func(c *uweb.Context) {
	 	 data := map[string]string {
	 	 	  "key": "value"
		 }		  	  
	 	 c.Render.Html(200, "account/login", data)
	 })	
	 
	 // post
	 uweb.Post("/api/login/", func(c *uweb.Context) {
	 	c.Render.Json(201, uweb.Map{
		  "key1": "value1",
		})
	 })
	 
	 // not support regexp match
	 uweb.Put("/account/:user_id", func (c *uweb.Context) {
	     userId := c.Req.Params["user_id"]
	 	 println(userId)
	 	 account.Noop(userId)
	 	 c.Res.Plain(201, "success")
     })
}

//
// src/models/account/noop.go
//
package account

func Noop(userId int) {
	// do nothing
}

```
