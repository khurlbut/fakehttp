package fakehttp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFakeHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Package 'github.com/khurlbut/fakehttp'")
}
