package anypointclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

/*
Copyright © 2021 Pontus Ullgren

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// AnypointClient is a API client for Anypoint Platform
type AnypointClient struct {
	HTTPClient *http.Client
	username   string
	password   string
	bearer     string
	baseURL    string
}

/*
NewAnypointClientWithToken creates a new Anypoint Client using the given token
*/
func NewAnypointClientWithToken(baseURL string, bearer string) AnypointClient {
	var c AnypointClient

	c.HTTPClient = &http.Client{}
	c.bearer = bearer
	c.baseURL = resolveBaseURLFromRegion(baseURL)
	return c
}

/*
NewAnypointClientWithCredentials creates a new Anypoint Client using the given username and password to aquire a token
*/
func NewAnypointClientWithCredentials(baseURL string, username string, password string) AnypointClient {
	var c AnypointClient

	c.HTTPClient = &http.Client{}
	c.baseURL = resolveBaseURLFromRegion(baseURL)
	c.username = username
	c.password = password
	c.bearer = c.getAuthorizationBearerToken()

	return c
}
func (client *AnypointClient) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s",
		client.baseURL,
		path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if client.bearer != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.bearer))
	}
	req.Header.Add("Cache-Control", "no-cache")
	return req, nil
}

func resolveBaseURLFromRegion(region string) string {
	switch strings.ToUpper(region) {
	case "EU":
		return "https://eu1.anypoint.mulesoft.com"
	case "US":
		return "https://anypoint.mulesoft.com"
	default:
		return region
	}
}
