package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/kashalls/unifi-client/internal/config"
	"golang.org/x/net/publicsuffix"
)

type Unifi struct {
	*config.Config
	*http.Client
	csrf string
}

func BuildClient(config *config.Config) (*Unifi, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &Unifi{
		Config: config,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipTLSVerify},
			},
			Jar: jar,
		},
	}

	if client.Config.ApiKey != "" {
		return client, nil
	}

	if err := client.Login(); err != nil {
		return nil, err
	}

	return client, nil
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func (c *Unifi) Login() error {
	body, err := json.Marshal(Login{
		Username: c.Config.Username,
		Password: c.Config.Password,
		Remember: c.Config.LongLogin,
	})
	if err != nil {
		return err
	}

	endpoint := LoginEndpoint
	if c.Config.ExternalController {
		endpoint = LoginExternalEndpoint
	}

	resp, err := c.doRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	if csrf := resp.Header.Get("x-csrf-token"); csrf != "" {
		c.csrf = resp.Header.Get("x-csrf-token")
	}
	return nil
}

func (c *Unifi) Logout() error {

	// Why are you trying to logout as an ApiKey?
	if c.Config.ApiKey != "" {
		return nil
	}

	_, err := c.doRequest("POST", LogoutEndpoint, nil)
	return err
}
