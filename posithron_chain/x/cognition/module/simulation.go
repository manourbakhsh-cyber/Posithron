package cognition

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/testutil/sample"
	cognitionsimulation "posithron/x/cognition/simulation"
	"posithron/x/cognition/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	cognitionGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		TaskRecordMap: []types.TaskRecord{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&cognitionGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateTask          = "op_weight_msg_cognition"
		defaultWeightMsgCreateTask int = 100
	)

	var weightMsgCreateTask int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTask, &weightMsgCreateTask, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTask = defaultWeightMsgCreateTask
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTask,
		cognitionsimulation.SimulateMsgCreateTask(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateTaskRecord          = "op_weight_msg_cognition"
		defaultWeightMsgCreateTaskRecord int = 100
	)

	var weightMsgCreateTaskRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTaskRecord, &weightMsgCreateTaskRecord, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTaskRecord = defaultWeightMsgCreateTaskRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTaskRecord,
		cognitionsimulation.SimulateMsgCreateTaskRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateTaskRecord          = "op_weight_msg_cognition"
		defaultWeightMsgUpdateTaskRecord int = 100
	)

	var weightMsgUpdateTaskRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTaskRecord, &weightMsgUpdateTaskRecord, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTaskRecord = defaultWeightMsgUpdateTaskRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTaskRecord,
		cognitionsimulation.SimulateMsgUpdateTaskRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteTaskRecord          = "op_weight_msg_cognition"
		defaultWeightMsgDeleteTaskRecord int = 100
	)

	var weightMsgDeleteTaskRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteTaskRecord, &weightMsgDeleteTaskRecord, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteTaskRecord = defaultWeightMsgDeleteTaskRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteTaskRecord,
		cognitionsimulation.SimulateMsgDeleteTaskRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
