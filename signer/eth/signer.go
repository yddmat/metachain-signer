package eth

// #cgo CFLAGS: -I../../wallet-core/include
// #cgo LDFLAGS: -L../../wallet-core/build -L../../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWAnySigner.h>
// #include <TrustWalletCore/TWEthereumChainID.h>
import "C"
import (
	"multichain-signer/signer/types"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Signer struct {
	privateKey []byte
}

func NewSigner(privateKey []byte) *Signer {
	return &Signer{privateKey: privateKey}
}

func (s *Signer) SignTx(txParams SignTxParams) ([]byte, error) {
	input := SigningInput{
		ChainId:    []byte{1},
		Nonce:      txParams.Nonce.Bytes(),
		GasPrice:   txParams.GasPrice.Bytes(),
		GasLimit:   txParams.GasLimit.Bytes(),
		ToAddress:  txParams.ToAddress,
		PrivateKey: s.privateKey,
		Transaction: &Transaction{TransactionOneof: &Transaction_Transfer_{
			Transfer: &Transaction_Transfer{
				Amount: txParams.Amount.Bytes(),
			},
		}},
	}

	inputBytes, err := proto.Marshal(&input)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling proto")
	}

	inputData := types.TWDataCreateWithGoBytes(inputBytes)
	defer C.TWDataDelete(inputData)

	outputData := C.TWAnySignerSign(inputData, C.TWCoinTypeEthereum)
	defer C.TWDataDelete(outputData)

	var output SigningOutput
	err = proto.Unmarshal(types.TWDataGoBytes(outputData), &output)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling sign")
	}

	return output.Encoded, nil
}
