package uweb

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	// secret in session
	CSRF_SECRET_KEY = "_csrf_secret"

	// token in session
	CSRF_TOKEN_KEY = "_csrf_token"
)

const (
	// the longer the better
	CSRF_SECRET_LEN = 18

	// doesn't need to be long
	CSRF_SALT_LEN = 8
)

// return true if except 
type CsrfExcepter func(*Context) bool

//
// CSRF middleware, depends on session
//
func MdCsrf(cookieAge int, f CsrfExcepter) Middleware {
	return NewCsrf(cookieAge, f)
}

//
// CSRF protect
//
type Csrf struct {
	cookieAge int
	exceptFunc CsrfExcepter
}

// Create csrf handler
func NewCsrf(cookieAge int, f CsrfExcepter) *Csrf {
	return &Csrf{
		cookieAge: cookieAge,
		exceptFunc: f,
	}
}

func (cf *Csrf) Name() string {
	return "csrf"
}

// Impl Middleware
func (cf *Csrf) Handle(c *Context) int {
	if cf.exceptFunc != nil {
		if cf.exceptFunc(c) {
			return NEXT_CONTINUE
		}
	}

	// lazily creates a csrf token
	// create one per session
	secret, token := c.Sess.Get(CSRF_SECRET_KEY), c.Sess.Get(CSRF_TOKEN_KEY)
	if len(secret) == 0 || len(token) == 0 {
		// create new token
		secret = cf.genSecret(CSRF_SECRET_LEN)
		salt := cf.genSalt(CSRF_SALT_LEN)
		token = cf.genToken(salt, secret)

		// save in session
		c.Sess.Set(CSRF_SECRET_KEY, secret)
		c.Sess.Set(CSRF_TOKEN_KEY, token)
		if DEBUG {
			log.Println(LOG_TAG, "Csrf: new token", token)
		}

		// for angular.js
		http.SetCookie(c.Res, &http.Cookie{
			Name:     "XSRF-TOKEN",
			Value:    token,
			Path:     "/",
			HttpOnly: false,
			MaxAge:   cf.cookieAge,
		})
	}

	// ignore method
	// DO NOT put this at beginning, as we need generating token at first request.
	if cf.bypassMethod(c.Req.Method) {
		return NEXT_CONTINUE
	}

	// parse reqToken
	reqToken := c.Req.FormValue("_csrf")
	if len(reqToken) == 0 {
		h := c.Req.Header
		reqToken = h.Get("X-CSRF-ReqToken")
		if len(reqToken) == 0 {
			reqToken = h.Get("X-XSRF-ReqToken")
		}
	}
	if len(reqToken) == 0 {
		c.Res.Status = 400
		c.Res.Err = errors.New("Csrf: no csrf")
		return NEXT_BREAK
	}

	// verify
	if err := cf.verify(secret, reqToken); err != nil {
		if DEBUG {
			log.Println("Csrf: verify error " + err.Error()) 
		}
		c.Res.Status = 403
		c.Res.Err = err
		return NEXT_BREAK
	}

	// ok
	return NEXT_CONTINUE
}

// bypass get-like methods
func (cf *Csrf) bypassMethod(m string) bool {
	switch m {
	case "GET", "HEAD", "OPTIONS":
		return true
	}
	return false
}

// create a secret key
// this __should__ be cryptographically secure,
// but generally client's can't/shouldn't-be-able-to access this so it really doesn't matt
func (cf *Csrf) genSecret(length int) string {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// create a random salt
func (cf *Csrf) genSalt(length int) string {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// create a csrf token
func (cf *Csrf) genToken(salt, secret string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, "-")
	io.WriteString(h, secret)
	return salt + "-" + base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (cf *Csrf) verify(secret, token string) error {
	// extract salt
	a := strings.SplitN(token, "-", 2)
	if len(a) != 2 {
		return errors.New("Csrf: invalid token 0")
	}
	salt := a[0]
	if len(salt) == 0 {
		return errors.New("Csrf: empty salt")
	}

	// token
	expected := cf.genToken(salt, secret)
	if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
		return errors.New("Csrf: invalid token 1")
	}

	// ok
	return nil
}
