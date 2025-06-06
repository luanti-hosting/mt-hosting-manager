package mtui

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (a *MtuiClient) DownloadZip(dir string) (io.ReadCloser, error) {
	req, err := a.request(http.MethodGet, "api/filebrowser/zip", nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	q := req.URL.Query()
	q.Set("dir", dir)
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do error: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (a *MtuiClient) AppendFile(filename string, offset int64, data []byte) error {
	req, err := a.request(http.MethodPut, "api/filebrowser/file", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}

	q := req.URL.Query()
	q.Set("filename", filename)
	q.Set("offset", fmt.Sprintf("%d", offset))
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	return nil
}

type AppendingWriter struct {
	filename string
	offset   int64
	cl       *MtuiClient
}

func (aw *AppendingWriter) Write(p []byte) (int, error) {
	err := aw.cl.AppendFile(aw.filename, aw.offset, p)
	if err != nil {
		return 0, err
	}
	aw.offset += int64(len(p))
	return len(p), nil
}

func (a *MtuiClient) UploadStream(filename string) io.Writer {
	return &AppendingWriter{
		filename: filename,
		cl:       a,
	}
}

func (a *MtuiClient) DeleteFile(filename string) error {
	req, err := a.request(http.MethodDelete, "api/filebrowser/file", nil)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}

	q := req.URL.Query()
	q.Set("filename", filename)
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	return nil
}

func (a *MtuiClient) UnzipFile(filename string) error {
	req, err := a.request(http.MethodPost, "api/filebrowser/unzip", nil)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}

	q := req.URL.Query()
	q.Set("filename", filename)
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	return nil
}

func (a *MtuiClient) GetDirectorySize(dir string) (int64, error) {
	req, err := a.request(http.MethodGet, "api/filebrowser/size", nil)
	if err != nil {
		return 0, fmt.Errorf("request error: %v", err)
	}

	q := req.URL.Query()
	q.Set("dir", dir)
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("readall error: %v", err)
	}

	size, err := strconv.ParseInt(string(resp_bytes), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parseint error: %v", err)
	}

	return size, nil
}
