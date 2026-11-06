package entanglement

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/testutil/sample"
	entanglementsimulation "posithron/x/entanglement/simulation"
	"posithron/x/entanglement/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	entanglementGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PairRecordMap: []types.PairRecord{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&entanglementGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreatePairRecord          = "op_weight_msg_entanglement"
		defaultWeightMsgCreatePairRecord int = 100
	)

	var weightMsgCreatePairRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePairRecord, &weightMsgCreatePairRecord, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePairRecord = defaultWeightMsgCreatePairRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePairRecord,
		entanglementsimulation.SimulateMsgCreatePairRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdatePairRecord          = "op_weight_msg_entanglement"
		defaultWeightMsgUpdatePairRecord int = 100
	)

	var weightMsgUpdatePairRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePairRecord, &weightMsgUpdatePairRecord, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePairRecord = defaultWeightMsgUpdatePairRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePairRecord,
		entanglementsimulation.SimulateMsgUpdatePairRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeletePairRecord          = "op_weight_msg_entanglement"
		defaultWeightMsgDeletePairRecord int = 100
	)

	var weightMsgDeletePairRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePairRecord, &weightMsgDeletePairRecord, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePairRecord = defaultWeightMsgDeletePairRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePairRecord,
		entanglementsimulation.SimulateMsgDeletePairRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgEntangleNodes          = "op_weight_msg_entanglement"
		defaultWeightMsgEntangleNodes int = 100
	)

	var weightMsgEntangleNodes int
	simState.AppParams.GetOrGenerate(opWeightMsgEntangleNodes, &weightMsgEntangleNodes, nil,
		func(_ *rand.Rand) {
			weightMsgEntangleNodes = defaultWeightMsgEntangleNodes
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgEntangleNodes,
		entanglementsimulation.SimulateMsgEntangleNodes(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
