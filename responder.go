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
	statusCode := fakeRequest.Response.StatusCode
	body := fakeRequest.Response.BodyBuffer
	responseHeader := fakeRequest.Response.Header

	if len(fakeRequest.Header) > 0 {
		err, s, b := validateHeaders(fakeRequest.Header, httpRequest.Header)
		if err != nil {
			statusCode = s
			body = []byte(b)
			responseHeader = make(http.Header)
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
	if (len(body)) > 0 {
		w.Write(body)
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
