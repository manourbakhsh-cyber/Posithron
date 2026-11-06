package entanglement

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"posithron/x/entanglement/types"
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
					RpcMethod: "ListPairRecord",
					Use:       "list-pair-record",
					Short:     "List all pair-record",
				},
				{
					RpcMethod:      "GetPairRecord",
					Use:            "get-pair-record [id]",
					Short:          "Gets a pair-record",
					Alias:          []string{"show-pair-record"},
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
					RpcMethod:      "CreatePairRecord",
					Use:            "create-pair-record [index] [ai-node] [energy-node] [status]",
					Short:          "Create a new pair-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "ai_node"}, {ProtoField: "energy_node"}, {ProtoField: "status"}},
				},
				{
					RpcMethod:      "UpdatePairRecord",
					Use:            "update-pair-record [index] [ai-node] [energy-node] [status]",
					Short:          "Update pair-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "ai_node"}, {ProtoField: "energy_node"}, {ProtoField: "status"}},
				},
				{
					RpcMethod:      "DeletePairRecord",
					Use:            "delete-pair-record [index]",
					Short:          "Delete pair-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "EntangleNodes",
					Use:            "entangle-nodes [ai-node] [energy-node]",
					Short:          "Send a entangle-nodes tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ai_node"}, {ProtoField: "energy_node"}},
				},
			},
		},
	}
}
