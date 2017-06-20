package uweb

import (
	"fmt"
)

//
// error pages
//
func MdErrPage(data Map) Middleware {
	return &errPage{
		data: data,
	}
}

//
// support 404
//
type errPage struct {
	data Map
}

func (e *errPage) Name() string {
	return "errors"
}

func (e *errPage) Handle(c *Context) int {
	if c.Req.Method != "GET" {
		return NEXT_CONTINUE
	}
	c.Next()
	
	if c.Res.Status >= 400 {
		if c.Res.Err != nil {
			e.data["error"] = c.Res.Err.Error()
			c.Res.Err = nil
		}
		c.Render.Html(fmt.Sprintf("errors/%d", c.Res.Status), e.data)
	}

	return NEXT_CONTINUE
}
