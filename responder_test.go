package fakehttp_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"

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
		fakeRequest = NewRequest()
		fakeRequest.Reply(200).BodyString("Body")
	})

	It("should write statusCode 200 when required headers are found", func() {
		httpRequest.Header.Set("requiredHeaderKey", "requiredHeaderValue")
		fakeRequest.SetHeader("requiredHeaderKey", "requiredHeaderValue")
		RequireHeadersResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockStatusCode).Should(Equal(200))
	})

	It("should write text in fakeRequest to the Body when required headers are found", func() {
		httpRequest.Header.Set("requiredHeaderKey", "requiredHeaderValue")
		fakeRequest.SetHeader("requiredHeaderKey", "requiredHeaderValue")
		RequireHeadersResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal("Body"))
	})

	It("should write statusCode 500 when required headers are not found", func() {
		fakeRequest.SetHeader("requiredHeaderKey", "requiredHeaderValue")
		RequireHeadersResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockStatusCode).Should(Equal(500))
	})

	It("should write missing headers to the Body when they are not found", func() {
		key := "Key"
		val := "Value"
		fakeRequest.SetHeader(key, val)
		RequireHeadersResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required header %s:%s not found!", key, val)))
	})

	It("should not mess up the request handler when headers are not found", func() {
		key := "Key"
		val := "Value"
		fakeRequest.SetHeader(key, val)
		RequireHeadersResponder(mockWriter, httpRequest, fakeRequest)
		Ω(mockHtmlBody).Should(Equal(fmt.Sprintf("500: Required header %s:%s not found!", key, val)))
		httpRequest.Header.Set(key, val)
		RequireHeadersResponder(mockWriter, httpRequest, fakeRequest)
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
