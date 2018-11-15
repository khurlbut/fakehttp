package fakehttp

import (
	"errors"
	"fmt"
	"net/http"
)

type Responder func(w http.ResponseWriter, r *http.Request, rh *Request)

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

func RequireHeadersResponder(w http.ResponseWriter, httpRequest *http.Request, fakeRequest *Request) {
	if len(fakeRequest.Header) > 0 {
		err, statusCode, body := validateHeaders(fakeRequest.Header, httpRequest.Header)
		if err != nil {
			fakeRequest.Response.StatusCode = statusCode
			fakeRequest.Response.BodyBuffer = []byte(body)
			fakeRequest.Response.Header = make(http.Header)
		}
	}
	if (len(fakeRequest.Response.Header)) > 0 {
		for k := range fakeRequest.Response.Header {
			w.Header().Add(k, fakeRequest.Response.Header.Get(k))
		}
	}
	if fakeRequest.Response.StatusCode > 0 {
		w.WriteHeader(fakeRequest.Response.StatusCode)
	}
	if (len(fakeRequest.Response.BodyBuffer)) > 0 {
		w.Write(fakeRequest.Response.BodyBuffer)
	}
}

func validateHeaders(requiredHeaders http.Header, incomingHeaders http.Header) (error, int, string) {
	for k, v := range requiredHeaders {
		if len(v) == 1 {
			requiredVal := v[0]
			val := incomingHeaders.Get(k)
			if val != requiredVal {
				fail := fmt.Sprintf("500: Required header %s:%s not found!", k, requiredVal)
				return errors.New("Fail"), 500, fail
			}
		}
	}
	return nil, 0, ""
}
