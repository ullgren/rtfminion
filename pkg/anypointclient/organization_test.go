package anypointclient

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Organisation", func() {
	It("should find the sub-org Master Organisation", func() {
		fixture, err := ioutil.ReadFile("testdata/login/successfull-login-response.json")
		if err != nil {
			Fail(fmt.Sprintf("Failed %v", err))
		}
		httpmock.RegisterResponder("POST", "/accounts/login", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(fixture))
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		})

		fixture2, err := ioutil.ReadFile("testdata/organization/me-response.json")
		httpmock.RegisterResponder("GET", "/accounts/api/me", func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(fixture2))
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		})

		token, err := client.Login()
		Ω(err == nil).Should(BeTrue(), "Error is %v", err)
		Ω(client.username).Should(Equal("user"), "user")
		Ω(token).Should(Equal("12345678-1234-1234-1234-123456789101"), "token")

		org, err := client.ResolveOrganisation("Example Inc/Example Inc Lab/bob-lab")
		Ω(err == nil).Should(BeTrue(), "Error is %v", err)
		Ω(org.ID).Should(Equal("12345678-0075-43d1-9395-aa55ffa06ea5"), "orgID")
	})

	It("should not call /accounts/api/me a second time", func() {
		org, err := client.ResolveOrganisation("Example Inc/Example Inc Lab/bob-lab")
		Ω(err == nil).Should(BeTrue(), "Error is %v", err)
		Ω(org.ID).Should(Equal("12345678-0075-43d1-9395-aa55ffa06ea5"), "orgID")
	})

	It("should also be able to find the master organisation (and still use the cache)", func() {
		org, err := client.ResolveOrganisation("Example Inc")
		Ω(err == nil).Should(BeTrue(), "Error is %v", err)
		Ω(org.ID).Should(Equal("12345678-6085-4179-9bed-917f6643df29"), "orgID")
	})
})
