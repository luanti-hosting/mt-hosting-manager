package btcscan

import "net/http"

type BTCScanClient struct {
	client http.Client
}

func New() *BTCScanClient {
	return &BTCScanClient{
		client: http.Client{},
	}
}
