package fakehttp

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	netURL "net/url"
	"strings"
)

// HTTPFake struct to hold a 'fake' server
type HTTPFake struct {
	server          *httptest.Server
	RequestHandlers []*Request
}

// Server build a new fake server
func Server() *HTTPFake {
	server := &HTTPFake{
		RequestHandlers: []*Request{},
	}

	server.server = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rh := server.findHandler(r)
		if rh == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("--- 404 Page Not Found"))
			return
		}
		if rh.CustomHandle != nil {
			rh.CustomHandle(w, r, rh)
			return
		}
		DefaultResponder(w, r, rh)
	}))

	return server
}

// Start the fake server
func (f *HTTPFake) Start(ip string, port string) *HTTPFake {
	f.server.Listener = listener(ip, port)
	f.server.Start()
	return f
}

// Close the fake server
func (f *HTTPFake) Close() {
	f.server.Close()
}

// URL of the fake server
func (f *HTTPFake) URL() string {
	return f.server.URL
}

func listener(ip string, port string) net.Listener {
	fmt.Printf("Attempting to listen on %s:%v\n", ip, port)
	l, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("--- TCP FAILED! Using TCP6! --- (err: " + err.Error())
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}

// NewHandler get a new request with a new handler
func (f *HTTPFake) NewHandler() *Request {
	rh := NewRequest()
	f.RequestHandlers = append(f.RequestHandlers, rh)
	return rh
}

// ResolveURL return the url used to reach the fake server
func (f *HTTPFake) ResolveURL(path string, args ...interface{}) string {
	format := f.server.URL + path
	return fmt.Sprintf(format, args...)
}

// Reset the fake server
func (f *HTTPFake) Reset() *HTTPFake {
	f.RequestHandlers = []*Request{}
	return f
}

func (f *HTTPFake) findHandler(r *http.Request) *Request {
	founds := []*Request{}
	url := r.URL.String()
	path := getURLPath(url)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	for _, rh := range f.RequestHandlers {
		if rh.Method != r.Method {
			continue
		}

		rhURL, _ := netURL.QueryUnescape(rh.URL.String())

		if rhURL == url {
			return rh
		}

		if strings.HasPrefix(rhURL, "*") {
			return rh
		}

		if getURLPath(rhURL) == path {
			founds = append(founds, rh)
		}
	}

	if len(founds) == 1 {
		return founds[0]
	}

	return nil
}

func getURLPath(url string) string {
	return strings.Split(url, "?")[0]
}
