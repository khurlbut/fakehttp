package fakehttp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"net/http"

	. "github.com/khurlbut/fakehttp"
)

var _ = Describe("HTTP Fake Tests", func() {
	var server *HTTPFake
	var ip = "127.0.0.1"
	var port = "8181"

	BeforeEach(func() {
		server = Server()
		server.Start(ip, port)
	})

	AfterEach(func() {
		server.Close()
	})

	It("should not be nil", func() {
		Ω(*server).ShouldNot(BeNil())
	})

	It("should intialize empty request handlers array", func() {
		Ω(server.RequestHandlers).ShouldNot(BeNil())
		Ω(len(server.RequestHandlers)).Should(BeZero())
	})

	It("should add a new Request to the array of Request Handlers", func() {
		r := server.NewHandler(false)
		Ω(len(server.RequestHandlers)).ShouldNot(BeZero())
		Ω(server.RequestHandlers[0]).Should(Equal(r))
	})

	It("should resolve the full URL to the server server for a given path", func() {
		resolvedURL := server.ResolveURL("%s?%s=%s", "/path/to/page", "param1", "value1")
		Ω(resolvedURL).Should(Equal("http://" + ip + ":" + port + "/path/to/page?param1=value1"))
	})

	It("should reset the Request Handler definitions", func() {
		server.NewHandler(false)
		server.Reset()
		Ω(len(server.RequestHandlers)).Should(BeZero())
	})

	It("should return the expected response on GET", func() {
		server.NewHandler(false).Get("/users").Reply(200).BodyString(`[{"username": "dreamer"}]`)

		res, err := http.Get(server.ResolveURL("/users"))
		Ω(err).ShouldNot(HaveOccurred())

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal(`[{"username": "dreamer"}]`))
	})

	It("should return the expected response for any path when * is used", func() {
		server.NewHandler(false).Get("*").Reply(200).BodyString("ALWAYS_THE_SAME")

		res, err := http.Get(server.ResolveURL("/"))
		Ω(err).ShouldNot(HaveOccurred())
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal("ALWAYS_THE_SAME"))

		res, err = http.Get(server.ResolveURL("/users"))
		Ω(err).ShouldNot(HaveOccurred())
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal("ALWAYS_THE_SAME"))

		res, err = http.Get(server.ResolveURL("/some/other/random/url"))
		Ω(err).ShouldNot(HaveOccurred())
		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal("ALWAYS_THE_SAME"))
	})

	It("should return 404", func() {
		res, err := http.Get(server.ResolveURL("/path/to/nowhere"))
		Ω(err).ShouldNot(HaveOccurred())

		defer res.Body.Close()
		Ω(res.StatusCode).Should(Equal(404))
	})

	It("should return 500 when using requires header handler without sending the headers", func() {
		expected := "500: Required header Key:value not found!\nHeaders --> map[User-Agent:[Go-http-client/1.1] Accept-Encoding:[gzip]]"
		fakeRequest := server.NewHandler(false).Get("/users").AddHeader("key", "value")
		fakeRequest.Reply(200).BodyString(`[{"username": "dreamer"}]`)
		fakeRequest.CustomHandle = SophisticatedResponder

		res, err := http.Get(server.ResolveURL("/users"))
		Ω(err).ShouldNot(HaveOccurred())

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		Ω(res.StatusCode).Should(Equal(500))
		Ω(string(body)).Should(Equal(expected))
	})

	It("should return properly when using requires header handler and sending the headers", func() {
		fakeRequest := server.NewHandler(false).Get("/users").AddHeader("key", "value")
		fakeRequest.Reply(200).BodyString(`[{"username": "dreamer"}]`)
		fakeRequest.CustomHandle = SophisticatedResponder

		client := &http.Client{}
		req, _ := http.NewRequest("GET", (server.ResolveURL("/users")), nil)
		req.Header.Add("key", "value")
		res, err := client.Do(req)
		Ω(err).ShouldNot(HaveOccurred())

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal(`[{"username": "dreamer"}]`))
	})

	It("should return 500 when using requires cookie handler without sending the cookie", func() {
		cookie := &http.Cookie{Name: "unknownShopperId", Value: "123"}
		fakeRequest := server.NewHandler(false).Get("/users").AddCookie(cookie)
		fakeRequest.Reply(200).BodyString(`[{"username": "dreamer"}]`)
		fakeRequest.CustomHandle = SophisticatedResponder

		res, err := http.Get(server.ResolveURL("/users"))
		Ω(err).ShouldNot(HaveOccurred())

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		Ω(res.StatusCode).Should(Equal(500))
		Ω(string(body)).Should(Equal("500: Required cookie unknownShopperId not found!"))
	})

	It("should return properly when using requires cooke handler and sending the cookie", func() {
		cookie := &http.Cookie{Name: "unknownShopperId", Value: "123"}

		fakeRequest := server.NewHandler(false).Get("/users").AddCookie(cookie)
		fakeRequest.Reply(200).BodyString(`[{"username": "dreamer"}]`)
		fakeRequest.CustomHandle = SophisticatedResponder

		client := &http.Client{}
		req, _ := http.NewRequest("GET", (server.ResolveURL("/users")), nil)
		req.AddCookie(cookie)
		res, err := client.Do(req)
		Ω(err).ShouldNot(HaveOccurred())

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal(`[{"username": "dreamer"}]`))

	})
})
