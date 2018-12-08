package fakehttp

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// Request - fake Request object
type Request struct {
	Method        string
	URL           *url.URL
	Response      *Response
	Header        http.Header
	CookieArray   []*http.Cookie
	CustomHandle  Responder
	InjectionKeys []string
}

// NewRequest - create a Request object
func NewRequest() *Request {
	return &Request{
		URL:      &url.URL{},
		Header:   make(http.Header),
		Response: NewResponse(),
	}
}

// Get - create a Get request object
func (r *Request) Get(path string) *Request {
	return r.method("GET", path)
}

// Post - create a Post request object
func (r *Request) Post(path string) *Request {
	return r.method("POST", path)
}

// Put - create a Put request object
func (r *Request) Put(path string) *Request {
	return r.method("PUT", path)
}

// Patch - create a Patch request object
func (r *Request) Patch(path string) *Request {
	return r.method("PATCH", path)
}

// Delete - create a Delete request object
func (r *Request) Delete(path string) *Request {
	return r.method("DELETE", path)
}

// Head - create a Head request object
func (r *Request) Head(path string) *Request {
	return r.method("HEAD", path)
}

// SetHeader - set a Header on a request object
func (r *Request) SetHeader(key string, val string) *Request {
	r.Header.Set(key, val)
	return r
}

// AddHeader - add a Header to a request object
func (r *Request) AddHeader(key string, val string) *Request {
	r.Header.Add(key, val)
	return r
}

// AddCookie - add a Cookie to a request object
func (r *Request) AddCookie(c *http.Cookie) *Request {
	r.CookieArray = append(r.CookieArray, c)
	return r
}

// Cookie - retrieve a cookie
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	for _, thisCookie := range r.CookieArray {
		if name == thisCookie.Name {
			return thisCookie, nil
		}
	}
	return nil, errors.New("cookie not found")
}

// Cookies - retrieve all the cookies
func (r *Request) Cookies() []*http.Cookie {
	return r.CookieArray
}

// Handle - add a Responder to a Request
func (r *Request) Handle(handle Responder) {
	r.CustomHandle = handle
}

// Reply - set the Status Code and return a Response
func (r *Request) Reply(status int) *Response {
	r.Response.Status(status)
	return r.Response
}

// AddInjectionKey - key used to inject a value from the http.Request into the Body of Response
func (r *Request) AddInjectionKey(key string) *Request {
	r.InjectionKeys = append(r.InjectionKeys, key)
	return r
}
func (r *Request) method(method, path string) *Request {
	r.URL.Path = normalize(path)
	r.Method = strings.ToUpper(method)
	return r
}

func normalize(p string) string {
	if strings.HasSuffix(p, "/") {
		return p
	}
	return p + "/"
}
