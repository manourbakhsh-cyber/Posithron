package greenenergy

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"posithron/x/greenenergy/types"
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
					RpcMethod: "ListEnergyRecord",
					Use:       "list-energy-record",
					Short:     "List all energy-record",
				},
				{
					RpcMethod:      "GetEnergyRecord",
					Use:            "get-energy-record [id]",
					Short:          "Gets a energy-record",
					Alias:          []string{"show-energy-record"},
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
					RpcMethod:      "CreateEnergyRecord",
					Use:            "create-energy-record [index] [amount] [source]",
					Short:          "Create a new energy-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "amount"}, {ProtoField: "source"}},
				},
				{
					RpcMethod:      "UpdateEnergyRecord",
					Use:            "update-energy-record [index] [amount] [source]",
					Short:          "Update energy-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "amount"}, {ProtoField: "source"}},
				},
				{
					RpcMethod:      "DeleteEnergyRecord",
					Use:            "delete-energy-record [index]",
					Short:          "Delete energy-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "SubmitEnergy",
					Use:            "submit-energy [amount] [source]",
					Short:          "Send a submit-energy tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "source"}},
				},
				{
					RpcMethod:      "RecordEnergy",
					Use:            "record-energy [amount] [source]",
					Short:          "Send a record-energy tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "source"}},
				},
			},
		},
	}
}
