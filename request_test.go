package fakehttp_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/khurlbut/fakehttp"
)

var _ = Describe("Request Tests", func() {
	var r *Request

	BeforeEach(func() {
		r = NewRequest(false)
	})

	It("should not create a Nil Request", func() {
		Ω(*r).ShouldNot(BeNil())
	})

	It("should initialize a URL pointer", func() {
		Ω(r.URL).ShouldNot(BeNil())
	})

	It("should initialize a Response pointer", func() {
		Ω(r.Response).ShouldNot(BeNil())
	})

	It("should mutate into a GET request with a path to GET", func() {
		r.Get("path/to/get")
		Ω(r.Method).Should(Equal("GET"))
		Ω(r.URL.Path).Should(Equal("path/to/get/"))
	})

	It("should mutate into a POST request with a path to POST", func() {
		r.Post("path/to/post")
		Ω(r.Method).Should(Equal("POST"))
		Ω(r.URL.Path).Should(Equal("path/to/post/"))
	})

	It("should mutate into a PUT request with a path to PUT", func() {
		r.Put("path/to/put")
		Ω(r.Method).Should(Equal("PUT"))
		Ω(r.URL.Path).Should(Equal("path/to/put/"))
	})

	It("should mutate into a PATCH request with a path to PATCH", func() {
		r.Patch("path/to/patch")
		Ω(r.Method).Should(Equal("PATCH"))
		Ω(r.URL.Path).Should(Equal("path/to/patch/"))
	})

	It("should mutate into a DELETE request with a path to DELETE", func() {
		r.Delete("path/to/delete")
		Ω(r.Method).Should(Equal("DELETE"))
		Ω(r.URL.Path).Should(Equal("path/to/delete/"))
	})

	It("should mutate into a HEAD request with a path to HEAD", func() {
		r.Head("path/to/head")
		Ω(r.Method).Should(Equal("HEAD"))
		Ω(r.URL.Path).Should(Equal("path/to/head/"))
	})

	It("should set a custom Responder onto the CustomHandler", func() {
		Ω(r.CustomHandle).Should(BeNil())
		r.Handle(DefaultResponder)
		Ω(r.CustomHandle).ShouldNot(BeNil())
	})

	Specify("that Reply will set a status on the contained Response", func() {
		res := r.Reply(200)
		Ω(res.StatusCode).Should(Equal(200))
	})

	It("should set a header value", func() {
		r.SetHeader("key", "value")
		Ω(r.Header.Get("key")).Should(Equal("value"))
	})

	It("should add a header value", func() {
		r.AddHeader("key", "value")
		Ω(r.Header.Get("key")).Should(Equal("value"))
	})

	It("should add and retrieve a cookie", func() {
		cookie := &http.Cookie{Name: "unknownShopperId", Value: "123"}
		r.AddCookie(cookie)
		cookie, err := r.Cookie("unknownShopperId")
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cookie.Value).Should(Equal("123"))
	})

	It("should add and retrieve multiple cookies", func() {
		cookie1 := &http.Cookie{Name: "cookie1", Value: "111"}
		cookie2 := &http.Cookie{Name: "cookie2", Value: "222"}
		cookie3 := &http.Cookie{Name: "cookie3", Value: "333"}

		r.AddCookie(cookie1)
		r.AddCookie(cookie2)
		r.AddCookie(cookie3)

		cookies := r.Cookies()

		Ω(cookies[0]).Should(Equal(cookie1))
		Ω(cookies[1]).Should(Equal(cookie2))
		Ω(cookies[2]).Should(Equal(cookie3))

		cookie, err := r.Cookie(cookie2.Name)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cookie.Value).Should(Equal("222"))
	})

	It("should add an InjectionKey", func() {
		r.AddInjectionKey("path")
		injectionKey := r.InjectionKeys[0]
		Ω(injectionKey).Should(Equal("path"))
	})

	It("should add a ServiceEndpoint", func() {
		r.AddServiceEndpoint("uri")
		serviceEndpoint := r.ServiceEndpoints[0]
		Ω(serviceEndpoint).Should(Equal("uri"))
	})
})
