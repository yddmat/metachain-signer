package eth

import (
	"encoding/json"
	"math/big"
)

type BigInt struct {
	big.Int
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	i.SetString(val, 10)

	return nil
}

type SignTxParams struct {
	ToAddress string `json:"to_address"`
	Amount    BigInt `json:"amount,string"`
	Nonce     BigInt `json:"nonce"`
	GasPrice  BigInt `json:"gas_price"`
	GasLimit  BigInt `json:"gas_limit"`
}
