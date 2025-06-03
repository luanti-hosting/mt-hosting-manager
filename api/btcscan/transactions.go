package btcscan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *BTCScanClient) GetAddressTransactions(address string) ([]*Transaction, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://btcscan.org/api/address/%s/txs", address), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	txs := []*Transaction{}
	err = json.NewDecoder(resp.Body).Decode(&txs)

	return txs, err
}
