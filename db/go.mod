module db

go 1.23.2

require (
	blockchain/chain v0.0.0-00010101000000-000000000000 // indirect
	blockchain/transaction v0.0.0-00010101000000-000000000000 // indirect
	blockchain/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
)

replace blockchain/chain => ../blockchain

replace blockchain/transaction => ../blockchain/transaction

replace blockchain/types => ../types
