package bitcoin

type UnspentTxOutput struct {
	Hash     string `json:"hash"`
	Index    uint32 `json:"index"`
	Sequence uint32 `json:"sequence"`
	Amount   int64  `json:"amount"`
}

type SignTxParams struct {
	UnspentTxOutput UnspentTxOutput `json:"utxo"`
	ToAddress       string          `json:"to_address"`
	ChangeAddress   string          `json:"change_address"`
	ByteFee         int64           `json:"byte_fee"`
	Amount          int64           `json:"amount"`
}
