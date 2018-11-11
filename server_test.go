package mockhttp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/khurlbut/mockhttp"
	"io/ioutil"
	"net/http"
)

var _ = Describe("HTTP Fake Tests", func() {
	var server *HTTPFake

	BeforeEach(func() {
		server = Server()
		server.Start()
	})

	AfterEach(func() {
		server.Server.Close()
	})

	It("should not be nil", func() {
		Ω(*server).ShouldNot(BeNil())
	})

	It("should intialize empty request handlers array", func() {
		Ω(server.RequestHandlers).ShouldNot(BeNil())
		Ω(len(server.RequestHandlers)).Should(BeZero())
	})

	It("should initialize Server", func() {
		Ω(server.Server).ShouldNot(BeNil())
	})

	It("should add a new Request to the array of Request Handlers", func() {
		r := server.NewHandler()
		Ω(len(server.RequestHandlers)).ShouldNot(BeZero())
		Ω(server.RequestHandlers[0]).Should(Equal(r))
	})

	/*
	 * This test demonstrates that the http.Server will generate a url pointing to
	 * localhost with a random port (5 digits).
	 *
	 * URL is: http://127.0.0.1:\d{5}/path/to/page?param1=value1
	 */
	It("should resolve the full URL to the server server for a given path", func() {
		resolvedURL := server.ResolveURL("%s?%s=%s", "/path/to/page", "param1", "value1")
		Ω(resolvedURL).Should(MatchRegexp("http:\\/\\/127\\.0\\.0\\.1:8181\\/path\\/to\\/page\\?param1=value1"))
		// Ω(resolvedURL).Should(MatchRegexp("http:\\/\\/127\\.0\\.0\\.1:\\d{5}\\/path\\/to\\/page\\?param1=value1"))
	})

	It("should reset the Request Handler definitions", func() {
		server.NewHandler()
		server.Reset()
		Ω(len(server.RequestHandlers)).Should(BeZero())
	})

	It("should return the expected response on GET", func() {
		server.NewHandler().Get("/users").Reply(200).BodyString(`[{"username": "dreamer"}]`)

		res, _ := http.Get(server.ResolveURL("/users"))
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		Ω(res.StatusCode).Should(Equal(200))
		Ω(string(body)).Should(Equal(`[{"username": "dreamer"}]`))

	})

	It("should return 404", func() {
		res, _ := http.Get(server.ResolveURL("/path/to/nowhere"))
		defer res.Body.Close()
		Ω(res.StatusCode).Should(Equal(404))
	})

	PIt("should just do stuff...", func() {
	})
})
