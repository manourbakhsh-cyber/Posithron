package cognition

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"posithron/x/cognition/types"
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
					RpcMethod: "ListTaskRecord",
					Use:       "list-task-record",
					Short:     "List all task-record",
				},
				{
					RpcMethod:      "GetTaskRecord",
					Use:            "get-task-record [id]",
					Short:          "Gets a task-record",
					Alias:          []string{"show-task-record"},
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
					RpcMethod:      "CreateTask",
					Use:            "create-task [prompt] [bounty]",
					Short:          "Send a create-task tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "prompt"}, {ProtoField: "bounty", Varargs: true}},
				},
				{
					RpcMethod:      "CreateTaskRecord",
					Use:            "create-task-record [index] [prompt] [bounty] [status]",
					Short:          "Create a new task-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "prompt"}, {ProtoField: "bounty"}, {ProtoField: "status"}},
				},
				{
					RpcMethod:      "UpdateTaskRecord",
					Use:            "update-task-record [index] [prompt] [bounty] [status]",
					Short:          "Update task-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "prompt"}, {ProtoField: "bounty"}, {ProtoField: "status"}},
				},
				{
					RpcMethod:      "DeleteTaskRecord",
					Use:            "delete-task-record [index]",
					Short:          "Delete task-record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
			},
		},
	}
}
