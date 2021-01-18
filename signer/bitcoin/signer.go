package bitcoin

// #cgo CFLAGS: -I../../wallet-core/include
// #cgo LDFLAGS: -L../../wallet-core/build -L../../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWBitcoinScript.h>
// #include <TrustWalletCore/TWAnySigner.h>
import "C"
import (
	"encoding/hex"
	"multichain-signer/signer/types"
	"unsafe"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Signer struct {
	privateKey []byte
	hashType   uint32

	addressPtr unsafe.Pointer
}

func NewSigner(privateKey []byte, hashType uint32, addressPtr unsafe.Pointer) *Signer {
	return &Signer{
		privateKey: privateKey,
		hashType:   hashType,
		addressPtr: addressPtr,
	}
}

func (s *Signer) SignTx(txParams SignTxParams) ([]byte, error) {
	utxoHash, err := hex.DecodeString(txParams.UnspentTxOutput.Hash)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid utxo hash")
	}

	script := C.TWBitcoinScriptLockScriptForAddress(s.addressPtr, C.TWCoinTypeBitcoin)
	defer C.TWBitcoinScriptDelete(script)

	scriptData := C.TWBitcoinScriptData(script)
	defer C.TWDataDelete(scriptData)

	input := SigningInput{
		HashType:      s.hashType,
		Amount:        txParams.Amount,
		ByteFee:       txParams.ByteFee,
		ToAddress:     txParams.ToAddress,
		ChangeAddress: txParams.ChangeAddress,
		PrivateKey:    [][]byte{s.privateKey},
		Utxo: []*UnspentTransaction{
			{
				OutPoint: &OutPoint{
					Hash:     utxoHash,
					Index:    txParams.UnspentTxOutput.Index,
					Sequence: txParams.UnspentTxOutput.Sequence,
				},
				Amount: txParams.UnspentTxOutput.Amount,
				Script: types.TWDataGoBytes(scriptData),
			},
		},
		CoinType: 0,
	}

	inputBytes, _ := proto.Marshal(&input)
	inputData := types.TWDataCreateWithGoBytes(inputBytes)
	defer C.TWDataDelete(inputData)

	outputData := C.TWAnySignerSign(inputData, C.TWCoinTypeBitcoin)
	defer C.TWDataDelete(outputData)

	var output SigningOutput
	err = proto.Unmarshal(types.TWDataGoBytes(outputData), &output)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling sign")
	}

	return output.Encoded, nil
}
