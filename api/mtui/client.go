package mtui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MtuiClient struct {
	client  http.Client
	url     string
	cookies []*http.Cookie
}

func New(url string) *MtuiClient {
	return &MtuiClient{
		url:     url,
		client:  http.Client{},
		cookies: make([]*http.Cookie, 0),
	}
}

func (a *MtuiClient) request(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", a.url, path), body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	for _, c := range a.cookies {
		req.AddCookie(c)
	}

	return req, nil
}

func (a *MtuiClient) Login(username, jwt_key string) error {
	req, err := a.request(http.MethodGet, fmt.Sprintf("api/loginadmin/%s", username), nil)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}

	q := req.URL.Query()
	q.Set("key", jwt_key)
	q.Set("disable_redirect", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	a.cookies = resp.Cookies()
	return nil
}

func (a *MtuiClient) GetStats() (*Stats, error) {
	req, err := a.request(http.MethodGet, "api/stats", nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	stats := &Stats{}
	err = json.NewDecoder(resp.Body).Decode(stats)
	if err != nil {
		return nil, fmt.Errorf("json parse error: %v", err)
	}

	return stats, nil
}
