package japicore

import (
	"strings"

	"github.com/JackalLabs/jackalgo/handlers/file_io_handler"
	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jutils"
)

func InitWalletSession() (*wallet_handler.WalletHandler, *file_io_handler.FileIoHandler) {
	seed := jutils.LoadEnvVarOrPanic("JAPI_SEED")
	rpc := jutils.LoadEnvVarOrFallback("JAPI_RPC", "https://jackal-testnet-rpc.polkachu.com:443")
	chainid := jutils.LoadEnvVarOrFallback("JAPI_CHAIN", "lupulella-2")
	operatingRoot := jutils.LoadEnvVarOrFallback("JAPI_OP_ROOT", "JAPI")

	if strings.HasPrefix(operatingRoot, "s/") {
		warning := "operatingRoot must not start with the \"s/\" prefix"
		panic(jutils.ProcessCustomError("InitWalletSession - HasPrefix", warning))
	}

	wallet, err := wallet_handler.NewWalletHandler(
		seed, // slim odor fiscal swallow piece tide naive river inform shell dune crunch canyon ten time universe orchard roast horn ritual siren cactus upon forum
		rpc,
		chainid)
	if err != nil {
		panic(err)
	}

	// fileIo, err := file_io_handler.NewFileIoHandler(wallet.WithGas("250000"))
	fileIo, err := file_io_handler.NewFileIoHandler(wallet)
	if err != nil {
		panic(err)
	}

	_, err = fileIo.DownloadFolder("s/" + operatingRoot)
	if err != nil {
		_, err := fileIo.GenerateInitialDirs([]string{operatingRoot})
		if err != nil {
			panic(err)
		}
	}

	return wallet, fileIo
}
