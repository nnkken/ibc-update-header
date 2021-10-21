package main

import (
	"fmt"

	"github.com/spf13/cobra"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	httprpcclient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/core"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/client/utils"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"

	ibctmtypes "github.com/cosmos/cosmos-sdk/x/ibc/light-clients/07-tendermint/types"
)

func newMarshaler() *codec.ProtoCodec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	std.RegisterInterfaces(interfaceRegistry)
	ibc.AppModuleBasic{}.RegisterInterfaces(interfaceRegistry)
	return marshaler
}

var Marshaler *codec.ProtoCodec = newMarshaler()

func getHeader(tmClient rpcclient.Client, height int64) (ibctmtypes.Header, int64) {
	clientCtx := client.Context{
		Client: tmClient,
		Height: height,
	}
	header, height, err := utils.QueryTendermintHeader(clientCtx)
	if err != nil {
		panic(err)
	}
	return header, height
}

func getIbcClientHeight(tmClient rpcclient.Client, clientID string) ibcclienttypes.Height {
	clientCtx := client.Context{
		Client: tmClient,
	}
	clientStateRes, err := utils.QueryClientState(clientCtx, clientID, false)
	if err != nil {
		panic(err)
	}
	var clientState exported.ClientState
	err = Marshaler.UnpackAny(clientStateRes.ClientState, &clientState)
	if err != nil {
		panic(err)
	}
	return clientState.GetLatestHeight().(ibcclienttypes.Height)
}

func main() {
	cmd := &cobra.Command{
		Use:   "ibc-update-header ibc-client-id host-chain-rpc remote-chain-rpc",
		Short: "Query IBC client state and headers to generate header JSON for `appd tx ibc tendermint-client update` command",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			ibcClientId := args[0]
			hostChainRpcEndpoint := args[1]
			remoteChainRpcEndpoint := args[2]

			hostTmClient, err := httprpcclient.New(hostChainRpcEndpoint, "/websocket")

			if err != nil {
				panic(err)
			}
			prevHeight := getIbcClientHeight(hostTmClient, ibcClientId)

			remoteTmClient, err := httprpcclient.New(remoteChainRpcEndpoint, "/websocket")
			if err != nil {
				panic(err)
			}

			oldHeader, _ := getHeader(remoteTmClient, int64(prevHeight.GetRevisionHeight()))
			latestHeader, _ := getHeader(remoteTmClient, 0)
			latestHeader.TrustedHeight = prevHeight
			latestHeader.TrustedValidators = oldHeader.ValidatorSet

			output := Marshaler.MustMarshalJSON(&latestHeader)
			fmt.Println(string(output))
		},
	}
	cmd.Execute()
}
