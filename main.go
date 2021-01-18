package main

//#cgo CFLAGS: -Iwallet-core/include
//#cgo LDFLAGS: -Lwallet-core/build -Lwallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
//#include <TrustWalletCore/TWHDWallet.h>
//#include <TrustWalletCore/TWPrivateKey.h>
//#include <TrustWalletCore/TWPublicKey.h>
//#include <TrustWalletCore/TWAnySigner.h>
import "C"
import (
	"fmt"
	"multichain-signer/api"
	"multichain-signer/signer"
	"multichain-signer/signer/bitcoin"
	"multichain-signer/signer/eth"
	"multichain-signer/signer/types"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Loading .env file:", err.Error())
		os.Exit(1)
	}

	str := types.TWStringCreateWithGoString(os.Getenv("SEED"))
	emtpy := types.TWStringCreateWithGoString("")
	defer C.TWStringDelete(str)
	defer C.TWStringDelete(emtpy)

	wallet := C.TWHDWalletCreateWithMnemonic(str, emtpy)
	defer C.TWHDWalletDelete(wallet)

	key := C.TWHDWalletGetKeyForCoin(wallet, C.TWCoinTypeBitcoin)
	keyData := C.TWPrivateKeyData(key)
	defer C.TWDataDelete(keyData)

	address := C.TWHDWalletGetAddressForCoin(wallet, C.TWCoinTypeBitcoin)
	defer C.TWStringDelete(address)

	server := api.Server{
		TxSigner: signer.NewSigner(
			bitcoin.NewSigner(types.TWDataGoBytes(keyData), 0, address),
			eth.NewSigner(types.TWDataGoBytes(keyData)),
		),

		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}

	fmt.Println("Starting server...")
	err := server.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
}
