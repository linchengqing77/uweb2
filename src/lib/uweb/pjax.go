package uweb

//
// Pjax middleware
//
func MdPjax() Middleware {
	return NewPjax()
}

//
// Pjax parser, only GET method.
//
type Pjax struct {
	// empty
}

// create pjax handler
func NewPjax() *Pjax {
	return new(Pjax)
}

// pjax
func (cf *Pjax) Name() string {
	return "pjax"
}

// Impl Middleware
func (pj *Pjax) Handle(c *Context) int {
	if c.Req.Method != "GET" {
		return NEXT_CONTINUE
	}
	pjax := false
	if c.Req.FormValue("_pjax") == "true" {
		pjax = true
	}
	if !pjax {
		h := c.Req.Header
		if h.Get("X-PJAX") == "true" {
			pjax = true
		}
	}
	c.Req.Pjax = pjax
	return NEXT_CONTINUE
}
