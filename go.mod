module github.com/nnkken/ibc-update-header

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.42.7
	github.com/spf13/cobra v1.2.1
	github.com/tendermint/tendermint v0.34.11
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
