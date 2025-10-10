package client

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// UniFi devices that have the controller built-in usually use this prefix.
	ControllerPrefix string = "/proxy/network"

	LoginEndpoint         string = "/api/auth/login"
	LoginExternalEndpoint string = "/api/login"
	LogoutEndpoint        string = "/api/logout"
)

func (c *Unifi) Endpoint(template string, params ...interface{}) (string, error) {
	host := c.Config.Host

	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = fmt.Sprint("http://", host)
	}

	u, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("invalid host URL: %w", err)
	}

	if !c.Config.ExternalController && !strings.HasPrefix(template, ControllerPrefix) && template != LoginEndpoint {
		template = fmt.Sprint(ControllerPrefix, template)
	}

	path := fmt.Sprint(template, params)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u = u.ResolveReference(&url.URL{Path: path})

	return u.String(), nil
}
