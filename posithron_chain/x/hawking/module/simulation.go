package hawking

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/testutil/sample"
	hawkingsimulation "posithron/x/hawking/simulation"
	"posithron/x/hawking/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	hawkingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		NodeRecordMap: []types.NodeRecord{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&hawkingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateNodeRecord          = "op_weight_msg_hawking"
		defaultWeightMsgCreateNodeRecord int = 100
	)

	var weightMsgCreateNodeRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateNodeRecord, &weightMsgCreateNodeRecord, nil,
		func(_ *rand.Rand) {
			weightMsgCreateNodeRecord = defaultWeightMsgCreateNodeRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateNodeRecord,
		hawkingsimulation.SimulateMsgCreateNodeRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateNodeRecord          = "op_weight_msg_hawking"
		defaultWeightMsgUpdateNodeRecord int = 100
	)

	var weightMsgUpdateNodeRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateNodeRecord, &weightMsgUpdateNodeRecord, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateNodeRecord = defaultWeightMsgUpdateNodeRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateNodeRecord,
		hawkingsimulation.SimulateMsgUpdateNodeRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteNodeRecord          = "op_weight_msg_hawking"
		defaultWeightMsgDeleteNodeRecord int = 100
	)

	var weightMsgDeleteNodeRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteNodeRecord, &weightMsgDeleteNodeRecord, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteNodeRecord = defaultWeightMsgDeleteNodeRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteNodeRecord,
		hawkingsimulation.SimulateMsgDeleteNodeRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgEvaporateStake          = "op_weight_msg_hawking"
		defaultWeightMsgEvaporateStake int = 100
	)

	var weightMsgEvaporateStake int
	simState.AppParams.GetOrGenerate(opWeightMsgEvaporateStake, &weightMsgEvaporateStake, nil,
		func(_ *rand.Rand) {
			weightMsgEvaporateStake = defaultWeightMsgEvaporateStake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgEvaporateStake,
		hawkingsimulation.SimulateMsgEvaporateStake(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
