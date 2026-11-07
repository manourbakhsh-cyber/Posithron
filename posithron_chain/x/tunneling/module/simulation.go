package tunneling

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/testutil/sample"
	tunnelingsimulation "posithron/x/tunneling/simulation"
	"posithron/x/tunneling/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tunnelingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		TrustRecordMap: []types.TrustRecord{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tunnelingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCheckTunnel          = "op_weight_msg_tunneling"
		defaultWeightMsgCheckTunnel int = 100
	)

	var weightMsgCheckTunnel int
	simState.AppParams.GetOrGenerate(opWeightMsgCheckTunnel, &weightMsgCheckTunnel, nil,
		func(_ *rand.Rand) {
			weightMsgCheckTunnel = defaultWeightMsgCheckTunnel
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCheckTunnel,
		tunnelingsimulation.SimulateMsgCheckTunnel(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateTrustRecord          = "op_weight_msg_tunneling"
		defaultWeightMsgCreateTrustRecord int = 100
	)

	var weightMsgCreateTrustRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTrustRecord, &weightMsgCreateTrustRecord, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTrustRecord = defaultWeightMsgCreateTrustRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTrustRecord,
		tunnelingsimulation.SimulateMsgCreateTrustRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateTrustRecord          = "op_weight_msg_tunneling"
		defaultWeightMsgUpdateTrustRecord int = 100
	)

	var weightMsgUpdateTrustRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTrustRecord, &weightMsgUpdateTrustRecord, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTrustRecord = defaultWeightMsgUpdateTrustRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTrustRecord,
		tunnelingsimulation.SimulateMsgUpdateTrustRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteTrustRecord          = "op_weight_msg_tunneling"
		defaultWeightMsgDeleteTrustRecord int = 100
	)

	var weightMsgDeleteTrustRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteTrustRecord, &weightMsgDeleteTrustRecord, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteTrustRecord = defaultWeightMsgDeleteTrustRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteTrustRecord,
		tunnelingsimulation.SimulateMsgDeleteTrustRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
