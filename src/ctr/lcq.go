package ctrl

import (
	"fmt"
	"lib/uweb"
	"strings"
)

func init() {
	uweb.Get("/index", Index)
	uweb.Get("/lcq", Lcq)
	uweb.Post("/login", Login)
}
func CsrfToken(c *uweb.Context) string {
	return c.Sess.Get(uweb.CSRF_TOKEN_KEY)
}

func Index(c *uweb.Context) (int, error) {
	return 200, c.Render.Html("ind", uweb.Map{
		"title":   "index page",
		"content": "welcome",
		"csrf":    CsrfToken(c),
	})
}
func Lcq(c *uweb.Context) (int, error) {
	return 200, c.Render.Html("lcq", uweb.Map{
		"title":   "hello lcq",
		"content": "my add first page",
		"csrf":    CsrfToken(c),
	})
}
func Login(c *uweb.Context) (int, error) {
	return 200, c.Render.Html("suc", uweb.Map{
		"title":   "login page",
		"content": "login success",
		"csrf":    CsrfToken(c),
	})
	//return 201, c.Res.Plain("write to new body")
}

// 过滤
type mdAuth struct{}

func MdAuth() uweb.Middleware {
	return new(mdAuth)
}

func (h *mdAuth) Name() string {
	return "" //request topName
}
func (h *mdAuth) Handle(c *uweb.Context) int {
	p := c.Req.URL.Path
	fmt.Println(p)

	// Name
	if strings.HasSuffix(p, "/lcq") || strings.HasSuffix(p, "/login") {
		return uweb.NEXT_CONTINUE
	}
	// root
	if p == "/" {
		c.Res.Redirect(p, "/index")
		return uweb.NEXT_BREAK
	}
	// 过滤条件
	if 1 == 1 && p != "/index" {
		//c.Res.Header().Set("xxxPage", "/xxx/xxx")
		c.Res.Redirect(p, "/index")
		return uweb.NEXT_BREAK
	}
	return uweb.NEXT_CONTINUE
}
