package tunneling

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"posithron/x/tunneling/types"
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
					RpcMethod: "ListTrustRecord",
					Use:       "list-trust-record",
					Short:     "List all trust-record",
				},
				{
					RpcMethod:      "GetTrustRecord",
					Use:            "get-trust-record [id]",
					Short:          "Gets a trust-record",
					Alias:          []string{"show-trust-record"},
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
					RpcMethod:      "CheckTunnel",
					Use:            "check-tunnel [trust-score] [required-score]",
					Short:          "Send a check-tunnel tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trust_score"}, {ProtoField: "required_score"}},
				},
				{
					RpcMethod:      "CreateTrustRecord",
					Use:            "create-trust-record [index] [score]",
					Short:          "Create a new trust-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "score"}},
				},
				{
					RpcMethod:      "UpdateTrustRecord",
					Use:            "update-trust-record [index] [score]",
					Short:          "Update trust-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "score"}},
				},
				{
					RpcMethod:      "DeleteTrustRecord",
					Use:            "delete-trust-record [index]",
					Short:          "Delete trust-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
			},
		},
	}
}
