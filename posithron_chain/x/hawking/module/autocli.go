package hawking

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"posithron/x/hawking/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListNodeRecord",
					Use:       "list-node-record",
					Short:     "List all node-record",
				},
				{
					RpcMethod:      "GetNodeRecord",
					Use:            "get-node-record [id]",
					Short:          "Gets a node-record",
					Alias:          []string{"show-node-record"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateNodeRecord",
					Use:            "create-node-record [index] [stake-amount]",
					Short:          "Create a new node-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "stake_amount"}},
				},
				{
					RpcMethod:      "UpdateNodeRecord",
					Use:            "update-node-record [index] [stake-amount]",
					Short:          "Update node-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "stake_amount"}},
				},
				{
					RpcMethod:      "DeleteNodeRecord",
					Use:            "delete-node-record [index]",
					Short:          "Delete node-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "EvaporateStake",
					Use:            "evaporate-stake [node-address] [network-total-stake]",
					Short:          "Send a evaporate-stake tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "node_address"}, {ProtoField: "network_total_stake"}},
				},
			},
		},
	}
}
