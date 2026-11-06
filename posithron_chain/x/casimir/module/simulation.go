package casimir

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/testutil/sample"
	casimirsimulation "posithron/x/casimir/simulation"
	"posithron/x/casimir/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	casimirGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PresenceRecordMap: []types.PresenceRecord{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&casimirGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreatePresenceRecord          = "op_weight_msg_casimir"
		defaultWeightMsgCreatePresenceRecord int = 100
	)

	var weightMsgCreatePresenceRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePresenceRecord, &weightMsgCreatePresenceRecord, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePresenceRecord = defaultWeightMsgCreatePresenceRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePresenceRecord,
		casimirsimulation.SimulateMsgCreatePresenceRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdatePresenceRecord          = "op_weight_msg_casimir"
		defaultWeightMsgUpdatePresenceRecord int = 100
	)

	var weightMsgUpdatePresenceRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePresenceRecord, &weightMsgUpdatePresenceRecord, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePresenceRecord = defaultWeightMsgUpdatePresenceRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePresenceRecord,
		casimirsimulation.SimulateMsgUpdatePresenceRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeletePresenceRecord          = "op_weight_msg_casimir"
		defaultWeightMsgDeletePresenceRecord int = 100
	)

	var weightMsgDeletePresenceRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePresenceRecord, &weightMsgDeletePresenceRecord, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePresenceRecord = defaultWeightMsgDeletePresenceRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePresenceRecord,
		casimirsimulation.SimulateMsgDeletePresenceRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgClaimCasimir          = "op_weight_msg_casimir"
		defaultWeightMsgClaimCasimir int = 100
	)

	var weightMsgClaimCasimir int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimCasimir, &weightMsgClaimCasimir, nil,
		func(_ *rand.Rand) {
			weightMsgClaimCasimir = defaultWeightMsgClaimCasimir
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimCasimir,
		casimirsimulation.SimulateMsgClaimCasimir(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
