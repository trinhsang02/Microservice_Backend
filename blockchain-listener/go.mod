module blockchain-listener

go 1.24.1

require (
	github.com/ethereum/go-ethereum v1.15.11
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgx/v5 v5.7.4
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)

require github.com/joho/godotenv v1.5.1

require github.com/yourusername/yourrepo/db v0.0.0

replace github.com/yourusername/yourrepo/db => ../db

require github.com/yourusername/yourrepo/mq v0.0.0

replace github.com/yourusername/yourrepo/mq => ../mq

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/bits-and-blooms/bitset v1.22.0 // indirect
	github.com/consensys/bavard v0.1.30 // indirect
	github.com/consensys/gnark-crypto v0.17.0 // indirect
	github.com/crate-crypto/go-eth-kzg v1.3.0 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/ethereum/c-kzg-4844/v2 v2.1.1 // indirect
	github.com/ethereum/go-verkle v0.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/supranational/blst v0.3.14 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
