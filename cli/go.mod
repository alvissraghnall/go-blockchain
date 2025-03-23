module cli

go 1.23.2

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/go-ethereum v1.14.12 // indirect
	github.com/holiman/uint256 v1.3.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	go-blockchain/wallet v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)

replace go-blockchain/wallet => ../wallet
