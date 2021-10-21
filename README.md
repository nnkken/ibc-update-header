# ibc-update-header

This is for generating header.json for Cosmos SDK based CLI tools to run the `appd tx ibc tendermint-client update [ibc-client-id] [header.json]` command.

This tool queries the latest client state (specifically the latest revision and height of the client), then query the remote chain's header and fill in the necessary fields for the IBC client update transaction.

## Build

`go build -o ibc-update-header main.go`

## Usage

`./ibc-update-header ibc-client-id host-chain-rpc remote-chain-rpc`

- `ibc-client-id` is the client ID existing on the chain
- `host-chain-rpc` is the RPC endpoint for the chain to run the IBC update client transaction
- `remote-chain-rpc` is the RPC endpoint for the chain which the client represents it's state

For example: `./ibc-update-header "07-tendermint-16" "http://some-likecoin-chain-rpc-endpoint:26657" "http://some-osmosis-chain-rpc-endpoint:26657"` outputs a JSON header to stdout for updating the IBC client existing on LikeCoin chain representing the state of the Osmosis chain.
