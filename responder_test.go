package fakehttp_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/khurlbut/fakehttp"
)

var _ = Describe("Responder Tests", func() {
	var fakeRequest *Request
	var mockWriter mockResponseWriter
	var httpRequest *http.Request

	BeforeEach(func() {
		mockStatusCode = -1
		mockHtmlBody = "empty"
		mockWriter = mockResponseWriter{header: make(http.Header)}
		httpRequest, _ = http.NewRequest("GET", "http://example.com", nil)
		fakeRequest = NewRequest(false)
		fakeRequest.Reply(200).BodyString("Body")
	})

	It("should write statusCode 200 when required headers are found", func() {
		httpRequest.Header.Set("requiredHeaderKey", "requiredHeaderValue")
		fakeRequest.SetHeader("requiredHeaderKey", "requiredHeaderValue")
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockStatusCode).Should(Equal(200))
	})

	It("should write text in fakeRequest to the Body when required headers are found", func() {
		httpRequest.Header.Set("requiredHeaderKey", "requiredHeaderValue")
		fakeRequest.SetHeader("requiredHeaderKey", "requiredHeaderValue")
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal("Body"))
	})

	It("should write statusCode 500 when required headers are not found", func() {
		fakeRequest.SetHeader("requiredHeaderKey", "requiredHeaderValue")
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockStatusCode).Should(Equal(500))
	})

	It("should write missing headers to the Body when they are not found", func() {
		key := "Key"
		val := "Value"
		fakeRequest.SetHeader(key, val)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required header %s:%s not found!", key, val)))
	})

	It("should not mess up the request handler when headers are not found", func() {
		key := "Key"
		val := "Value"
		fakeRequest.SetHeader(key, val)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required header %s:%s not found!", key, val)))
		httpRequest.Header.Set(key, val)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal("Body"))
	})

	It("should write incoming headers to the Body when required headers are not found", func() {
		httpRequest.Header.Set("requiredHeaderKey", "requiredHeaderValue")
		key := "Key"
		val := "Value"
		fakeRequest.SetHeader(key, val)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required header %s:%s not found!\nHeaders --> map[Requiredheaderkey:[requiredHeaderValue]]", key, val)))
	})

	It("should write statusCode 200 when required cookie is found", func() {
		cookieInServer := &http.Cookie{Name: "cookie", Value: "111"}
		cookieInRequest := &http.Cookie{Name: "cookie", Value: "111"}
		httpRequest.AddCookie(cookieInServer)
		fakeRequest.AddCookie(cookieInRequest)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockStatusCode).Should(Equal(200))
	})

	It("should write text in fakeRequest to the Body when required cookie is found", func() {
		cookieInServer := &http.Cookie{Name: "cookie", Value: "111"}
		cookieInRequest := &http.Cookie{Name: "cookie", Value: "111"}
		httpRequest.AddCookie(cookieInServer)
		fakeRequest.AddCookie(cookieInRequest)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal("Body"))
	})

	It("should write statusCode 500 when required cookie is not found", func() {
		cookieInServer := &http.Cookie{Name: "cookie", Value: "111"}
		fakeRequest.AddCookie(cookieInServer)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockStatusCode).Should(Equal(500))
	})

	It("should write missing cookie to the Body when it is not found", func() {
		cookieInServer := &http.Cookie{Name: "cookie", Value: "111"}
		fakeRequest.AddCookie(cookieInServer)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required cookie %s not found!", cookieInServer.Name)))
	})

	It("should not mess up the request handler when cookie is not found", func() {
		cookieInServer := &http.Cookie{Name: "cookie", Value: "111"}
		cookieInRequest := &http.Cookie{Name: "cookie", Value: "111"}
		fakeRequest.AddCookie(cookieInRequest)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required cookie %s not found!", cookieInRequest.Name)))
		httpRequest.AddCookie(cookieInServer)
		SophisticatedResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal("Body"))
	})

})

var mockStatusCode int
var mockHtmlBody string

type mockResponseWriter struct {
	header http.Header
}

func (w mockResponseWriter) Header() http.Header {
	return w.header
}

func (w mockResponseWriter) Write(bytes []byte) (int, error) {
	mockHtmlBody = string(bytes)
	return len(mockHtmlBody), nil
}

func (w mockResponseWriter) WriteHeader(statusCode int) {
	mockStatusCode = statusCode
}
