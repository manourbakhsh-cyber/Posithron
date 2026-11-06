package casimir

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"posithron/x/casimir/types"
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
					RpcMethod: "ListPresenceRecord",
					Use:       "list-presence-record",
					Short:     "List all presence-record",
				},
				{
					RpcMethod:      "GetPresenceRecord",
					Use:            "get-presence-record [id]",
					Short:          "Gets a presence-record",
					Alias:          []string{"show-presence-record"},
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
					RpcMethod:      "CreatePresenceRecord",
					Use:            "create-presence-record [index] [status] [last-seen]",
					Short:          "Create a new presence-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "status"}, {ProtoField: "last_seen"}},
				},
				{
					RpcMethod:      "UpdatePresenceRecord",
					Use:            "update-presence-record [index] [status] [last-seen]",
					Short:          "Update presence-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "status"}, {ProtoField: "last_seen"}},
				},
				{
					RpcMethod:      "DeletePresenceRecord",
					Use:            "delete-presence-record [index]",
					Short:          "Delete presence-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "ClaimCasimir",
					Use:            "claim-casimir [node-address]",
					Short:          "Send a claim-casimir tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "node_address"}},
				},
			},
		},
	}
}
