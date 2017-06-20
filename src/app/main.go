package main

import (
	"ctr"
	"lib/uweb"
)

func main() {
	// uweb
	uweb.DEBUG = true
	uweb.DEVELOPMENT = true
	uweb.SID_COOKIE_KEY = "_uweb_sid"
	// app
	app := uweb.NewApp()
	// hacheck
	app.Use(uweb.MdIgnore([]string{"/hacheck"}))
	// static
	app.Use(uweb.MdFavicon("../pub/img/favicon.ico"))
	app.Use(uweb.MdStatic("/public", "../pub")) // before compress
	// compress
	app.Use(uweb.MdCompress())
	// log
	app.Use(uweb.MdLogger(uweb.LOG_LEVEL_2))
	// session
	app.Use(uweb.MdCache("memcache", "localhost:11211"))
	app.Use(uweb.MdSession(3600 * 24 * 29))
	app.Use(uweb.MdFlash())
	// csrf
	app.Use(uweb.MdCsrf(3600*24*28, func(c *uweb.Context) bool {
		return false
	}))
	// render
	app.Use(uweb.MdRender("../pub/", ".html", "[[", "]]"))
	// error page
	app.Use(uweb.MdErrPage(uweb.Map{
		"404_home_url": "",
	}))
	// @ctrl/auth.go
	app.Use(ctrl.MdAuth())
	// router
	app.Use(uweb.MdRouter())
	// listen
	app.Listen(":8888")
}
