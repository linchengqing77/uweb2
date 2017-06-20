package uweb

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"fmt"
	"io"
	"path"
	"os"
	
	"lib/uuid"
)

//
// Params
//
type Params map[string]string

// Convert to int value, if fail return 0
func (p Params) Int64(key string) int64 {
	s, ok := p[key]
	if !ok {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(LOG_TAG, "Params: Int64 err", err)
		return 0
	}
	return v
}

// -----------------------------------------------------------------------------
// Request

//
// Wrap http request
//
type Request struct {
	// embbed request for convenient
	*http.Request

	// client ip
	IP string

	// url pattern params, Router middleware will set it
	Params Params

	// should work with pjax middleware
	Pjax bool
}

// Create request
func NewRequest(req *http.Request) *Request {
	return &Request{req, readIp(req), nil, false}
}

// parse real ip if possible
// may use nginx to set header fields
func readIp(r *http.Request) string {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		i := strings.Index(v, ", ")
		if i == -1 {
			i = len(v)
		}
		return v[:i]
	}
	if v := r.Header.Get("X-Real-IP"); v != "" {
		return v
	}
	return ""
}

// -----------------------------------------------------------------------------
// url

// get full url for p
func (r *Request) UrlFor(p string) string {
	scheme := "http://"
    if r.TLS != nil {
        scheme = "https://"
    }
	return fmt.Sprintf("%s%s%s", scheme, r.Host, p)
}

// get current full url
func (r *Request) UrlCurr() string {
	return r.UrlFor(r.RequestURI)
}

// -----------------------------------------------------------------------------
// form utils

func (r *Request) FormStr(k string, min, max int) (string, error) {
	v := r.FormValue(k)
	if len(v) < min || len(v) > max {
		return "", fmt.Errorf("len(%s) should range in:[%d, %d] ", k, min, max)
	}
	return v, nil
}

func (r *Request) FormInt64(k string) (int64, error) {
	v := r.FormValue(k)
	if len(v) == 0 {
		return 0, fmt.Errorf("%s is empty", k)
	}
	return strconv.ParseInt(v, 10, 64)
}

func (r *Request) FormUint64(k string) (uint64, error) {
	v := r.FormValue(k)
	if len(v) == 0 {
		return 0, fmt.Errorf("%s is empty", k)
	}
	return strconv.ParseUint(v, 10, 64)
}

func (r *Request) FormInt(k string) (int, error) {
	i64, err := r.FormInt64(k)
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

func (r *Request) FormFloat64(k string) (float64, error) {
	v := r.FormValue(k)
	if len(v) == 0 {
		return 0, fmt.Errorf("%s is empty", k)
	}
	return strconv.ParseFloat(v, 64)
}

func (r *Request) FormStrArray(k string) []string {
	r.ParseForm()
	return r.Form[k]
}

func (r *Request) FormInt64Array(k string) ([]int64, error) {
	srr := r.FormStrArray(k)
	if len(srr) == 0 {
		return nil, nil
	}
	
	var irr []int64
	for _, s := range srr {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		irr = append(irr, i)
	}
	
	return irr, nil
}

//
// 示例：c.Req.FormFileSave("fire_cert_file", etc.Path("/data/upload"))
//
func (r *Request) FormFileSave(name, dstDir string) (string, error) {
	// 主要为了调用ParseMulitipartForm，记得先关闭这个src
	if src, _, err := r.FormFile(name); err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	} else {
		src.Close()
	}

	//
	// h5可以上传多个文件了，我们把所有文件名用逗号拼接返回:
	// aaa.jpg,bbb.jpg
	//
	dstAry := []string{}
	fhs := r.MultipartForm.File[name]
	for i, _ := range fhs {
		dstName, err := func(i int) (string, error) {
			// f
			f := fhs[i]

			// src
			src, err := f.Open()
			if err != nil {
				return "", err
			}
			defer src.Close()

			// dst
			dstName := uuid.New() + path.Ext(f.Filename)
			dst, err := os.Create(path.Join(dstDir, dstName))
			if err != nil {
				return "", err
			}
			defer dst.Close()

			// copy
			if _, err := io.Copy(dst, src); err != nil {
				return "", err
			}

			// ok
			return dstName, nil
		}(i)

		if err != nil {
			return "", err
		}
		dstAry = append(dstAry, dstName)
	}

	return strings.Join(dstAry, ","), nil
}