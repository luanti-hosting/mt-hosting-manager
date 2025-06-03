package btcscan

type Transaction struct {
	TxID    string              `json:"txid"`
	Version int                 `json:"version"`
	Vout    []*TransactionEntry `json:"vout"`
	Status  TransactionStatus   `json:"status"`
}

type TransactionEntry struct {
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

type TransactionStatus struct {
	Confirmed bool `json:"confirmed"`
}
