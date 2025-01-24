module transaction

go 1.23.2

replace blockchain/wallet => ../wallet
replace blockchain/chain => ../blockchain

require (
	blockchain/wallet v0.0.0-00010101000000-000000000000 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/go-ethereum v1.14.12 // indirect
	github.com/holiman/uint256 v1.3.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)
