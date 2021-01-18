package signer

import (
	"encoding/json"
	"multichain-signer/signer/bitcoin"
	"multichain-signer/signer/eth"

	"github.com/pkg/errors"
)

const (
	TypeBitcoin = "bitcoin"
	TypeEth     = "eth"
)

var (
	ErrUnknownSignerType = errors.New("unknown signer type")
	ErrInvalidTxData     = errors.New("invalid tx params data")
)

type signer struct {
	bitcoinSigner *bitcoin.Signer
	ethSigner     *eth.Signer
}

func NewSigner(bitcoinSigner *bitcoin.Signer, ethSigner *eth.Signer) *signer {
	return &signer{bitcoinSigner: bitcoinSigner, ethSigner: ethSigner}
}

func (s *signer) Sign(signerType string, txData json.RawMessage) ([]byte, error) {
	switch signerType {
	case TypeEth:
		params := eth.SignTxParams{}
		err := json.Unmarshal(txData, &params)
		if err != nil {
			return nil, ErrInvalidTxData
		}

		sig, err := s.ethSigner.SignTx(params)
		if err != nil {
			return nil, errors.Wrapf(err, "signing eth tx")
		}

		return sig, nil
	case TypeBitcoin:
		params := bitcoin.SignTxParams{}
		err := json.Unmarshal(txData, &params)
		if err != nil {
			return nil, ErrInvalidTxData
		}

		sig, err := s.bitcoinSigner.SignTx(params)
		if err != nil {
			return nil, errors.Wrapf(err, "signing bitcoin tx")
		}

		return sig, nil
	default:
		return nil, ErrUnknownSignerType
	}
}
