package mtui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *MtuiClient) CreateRestoreJob(job *CreateRestoreJob) (*RestoreJobInfo, error) {
	obj, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("json error: %v", err)
	}

	req, err := a.request(http.MethodPost, "api/restorejob", bytes.NewReader(obj))
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

	info := &RestoreJobInfo{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, fmt.Errorf("json response error: %v", err)
	}

	return info, nil
}

func (a *MtuiClient) GetRestoreJobInfo(id string) (*RestoreJobInfo, error) {
	req, err := a.request(http.MethodGet, fmt.Sprintf("api/restorejob/%s", id), nil)
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

	info := &RestoreJobInfo{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, fmt.Errorf("json response error: %v", err)
	}

	return info, nil
}
