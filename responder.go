package fakehttp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Responder - an abstract response builder
type Responder func(w http.ResponseWriter, r *http.Request, rh *Request)

// DefaultResponder - the default (simple) responder
func DefaultResponder(w http.ResponseWriter, r *http.Request, rh *Request) {
	if (len(rh.Response.Header)) > 0 {
		for k := range rh.Response.Header {
			w.Header().Add(k, rh.Response.Header.Get(k))
		}
	}
	if rh.Response.StatusCode > 0 {
		w.WriteHeader(rh.Response.StatusCode)
	}
	if (len(rh.Response.BodyBuffer)) > 0 {
		w.Write(rh.Response.BodyBuffer)
	}
}

// SophisticatedResponder - responder with validation of headers and cookies, body insertion of data from request and invocation of service dependencies
func SophisticatedResponder(w http.ResponseWriter, httpRequest *http.Request, fakeRequest *Request) {
	statusCode := fakeRequest.Response.StatusCode
	body := fakeRequest.Response.BodyBuffer
	responseHeader := fakeRequest.Response.Header

	if len(fakeRequest.Header) > 0 {
		s, b, err := validateHeaders(fakeRequest.Header, httpRequest.Header)
		if err != nil {
			statusCode = s
			body = []byte(b)
		}
	}
	if len(fakeRequest.Cookies()) > 0 {
		s, b, err := validateCookies(fakeRequest.Cookies(), httpRequest.Cookies())
		if err != nil {
			statusCode = s
			body = []byte(b)
		}
	}
	if (len(responseHeader)) > 0 {
		for k := range fakeRequest.Response.Header {
			w.Header().Add(k, responseHeader.Get(k))
		}
	}
	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}
	serviceResponses := ""
	if len(fakeRequest.ServiceEndpoints) > 0 {
		for _, uri := range fakeRequest.ServiceEndpoints {
			status, body, err := invokeServiceEndpoint(uri)
			if err == nil {
				serviceResponses += (uri + ": ")
				serviceResponses += (status + ": ")
				serviceResponses += (body)
				serviceResponses += "<br>"
			}
		}
	}
	if (len(body)) > 0 {
		b := string(body)
		if len(fakeRequest.InjectionKeys) > 0 {
			for _, k := range fakeRequest.InjectionKeys {
				if k == "path" {
					body = []byte(fmt.Sprintf(b, strings.TrimPrefix(httpRequest.URL.Path, "/")))
					b = string(body)
				}
			}
		}

		if len(serviceResponses) > 0 {
			b += serviceResponses
		}

		if fakeRequest.RenderHTML {
			b = "<html><head><title>fakeserver</title></head><body>" + b + "</body></html>"
		}

		w.Write([]byte(b))
	}
}

func validateHeaders(requiredHeaders http.Header, incomingHeaders http.Header) (int, string, error) {
	for k, v := range requiredHeaders {
		if len(v) == 1 {
			requiredVal := v[0]
			val := incomingHeaders.Get(k)
			if val != requiredVal {
				fail := fmt.Sprintf("500: Required header %s:%s not found!", k, requiredVal)
				if len(incomingHeaders) > 0 {
					fail = fail + fmt.Sprintf("\nHeaders --> %v", incomingHeaders)
				}
				return 500, fail, errors.New("Fail")
			}
		}
	}
	return 0, "", nil
}

func validateCookies(requiredCookies []*http.Cookie, incomingCookies []*http.Cookie) (int, string, error) {
	for _, cookie := range requiredCookies {
		if nil == findCookie(cookie.Name, incomingCookies) {
			fail := fmt.Sprintf("500: Required cookie %s not found!", cookie.Name)
			return 500, fail, errors.New("Fail")
		}
	}
	return 0, "", nil
}

func findCookie(name string, cookieArray []*http.Cookie) *http.Cookie {
	for _, thisCookie := range cookieArray {
		if name == thisCookie.Name {
			return thisCookie
		}
	}
	return nil
}

func invokeServiceEndpoint(uri string) (string, string, error) {
	response, err := http.Get(uri)
	if err != nil {
		log.Printf("Error invoking service endpoint %s: %v", uri, err)
		return "500", "", err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("%s", err)
		return "500", "", err
	}

	return response.Status, string(contents), nil
}
