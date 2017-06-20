package uweb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

//
// Http response
//
type Response struct {
	http.ResponseWriter

	Status int
	Err    error
	Body   []byte

	Close func()
}

// Create response with response
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w, 0, nil, nil, nil}
}

// Send status and body
func (res *Response) End(req *Request) error {
	// if error, ignore others
	if res.Err != nil {
		http.Error(res, res.Err.Error(), res.Status)
		return nil
	}

	// fix status
	if res.Status == 0 {
		switch req.Method {
		case "GET":
			res.Status = http.StatusOK
		case "POST", "PUT":
			res.Status = http.StatusCreated
		case "DELETE":
			res.Status = http.StatusNoContent
		default:
			res.Status = http.StatusOK
		}
		if res.Err != nil {
			res.Status = http.StatusInternalServerError
		}
	}

	// fix content-xxx
	if len(res.Body) > 0 {
		if ct := res.Header().Get("Content-Type"); len(ct) == 0 {
			res.Header().Set("Content-Type", http.DetectContentType(res.Body))
		}
	} else {
		res.Status = 204
		res.Header().Del("Content-Type")
		res.Header().Del("Content-Length")
		res.Header().Del("Content-Encoding")
	}

	// write body
	res.WriteHeader(res.Status)
	if _, err := res.Write(res.Body); err != nil {
		return err
	}

	// release if needed
	if res.Close != nil {
		res.Close()
	}

	// ok
	return nil
}

// Plain text
func (res *Response) Plain(data string) error {
	w := res
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Body = []byte(data)
	return nil
}

// about jsonp see:
// http://www.cnblogs.com/dowinning/archive/2012/04/19/json-jsonp-jquery.html
func (res *Response) Jsonp(padding string, v interface{}) error {
	// w
	w := res

	// body
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if len(padding) > 0 {
		result = []byte(fmt.Sprintf("%s(%s);", padding, string(result)))
	}
	w.Body = result

	// header
	w.Header().Del("Content-Length")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// ok
	return nil
}

// json
func (res *Response) Json(v interface{}) error {
	return res.Jsonp("", v)
}

// Html
func (res *Response) Html(body []byte) error {
	w := res
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Body = body
	return nil
}

// from - c.Req.URL.Path
// to - new path
func (res *Response) Redirect(from, to string) error {
	// copy from http/server.go
	// Location should be an absolute URI, like
	if u, err := url.Parse(to); err == nil {
		oldpath := from
		if oldpath == "" { // should not happen, but avoid a crash if it does
			oldpath = "/"
		}
		if u.Scheme == "" {
			// no leading http://server
			if to == "" || to[0] != '/' {
				// make relative path absolute
				olddir, _ := path.Split(oldpath)
				to = olddir + to
			}
			var query string
			if i := strings.Index(to, "?"); i != -1 {
				to, query = to[:i], to[i:]
			}
			// clean up but preserve trailing slash
			trailing := strings.HasSuffix(to, "/")
			to = path.Clean(to)
			if trailing && !strings.HasSuffix(to, "/") {
				to += "/"
			}
			to += query
		}
	}

	// RFC2616 recommends that a short note "SHOULD" be included in the
	// response because older user agents may not understand 301/307.
	// Shouldn't send the response for POST or HEAD; that leaves GET.
	res.Status = 302
	res.Header().Set("Location", to)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Body = []byte("Redirecting to " + to + ".")
	
	// ok
	return nil
}