package testutils

import (
	"fmt"
	"net/http"
)

const (
	MockToken          = "this_is_cc_token"
	MockRefreshToken = "this_is_refresh_token"
	ssoToken         = "this_is_sso_token"
)

type MockCommander struct {
	url string
}

func (c *MockCommander) Start() error {
	// Use a go routine so this function immediately returns to caller to
	// continue processing (needed for SSO flow)
	go func() {
		url := fmt.Sprintf("%s?code=%s", c.url, ssoToken)
		_, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error calling redirect: %s\n", url)
			return
		}
	}()

	return nil
}

func (c *MockCommander) GetURL(url string) {
	c.url = url
	return
}
