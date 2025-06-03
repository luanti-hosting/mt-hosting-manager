package btcscan_test

import (
	"mt-hosting-manager/api/btcscan"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	c := btcscan.New()
	txid := "bc1qwaelfcdg5v7a4v6k785h40cd2y66pws9tjqdjg"
	txs, err := c.GetAddressTransactions(txid)
	assert.NoError(t, err)
	assert.NotNil(t, txs)
	assert.True(t, len(txs) > 0)
	assert.NotNil(t, txs[0].Status)
}
