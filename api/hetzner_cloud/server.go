package hetzner_cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *HetznerCloudClient) CreateServer(csr *CreateServerRequest) (*CreateServerResponse, error) {
	data, err := json.Marshal(csr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.hetzner.cloud/v1/servers", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	crsresp := &CreateServerResponse{}
	err = json.Unmarshal(resp_bytes, crsresp)
	return crsresp, err
}

func (c *HetznerCloudClient) GetServer(id string) (*CreateServerResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.hetzner.cloud/v1/servers/%s", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	csr := &CreateServerResponse{}
	err = json.NewDecoder(resp.Body).Decode(csr)

	return csr, err
}

var ErrStillInPowerCycle = errors.New("server still in previous power-cycle")

func (c *HetznerCloudClient) PowerOffServer(id string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.hetzner.cloud/v1/servers/%s/actions/poweroff", id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 423 {
		return ErrStillInPowerCycle
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HetznerCloudClient) PowerOnServer(id string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.hetzner.cloud/v1/servers/%s/actions/poweron", id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 423 {
		return ErrStillInPowerCycle
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HetznerCloudClient) DeleteServer(id string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://api.hetzner.cloud/v1/servers/%s", id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return nil
}
