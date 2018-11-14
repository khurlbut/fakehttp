package fakehttp

import (
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
	keys := len(fakeRequest.Header)
	if keys > 0 {
		for k, v := range fakeRequest.Header {
			if len(v) == 1 {
				requiredVal := v[0]
				val := httpRequest.Header.Get(k)
				if val != requiredVal {
					w.Write([]byte(fmt.Sprintf("500: Required header %s:%s not found!", k, requiredVal)))
					w.WriteHeader(500)
					return
				}
			}
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
