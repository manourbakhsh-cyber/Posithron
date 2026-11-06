package heisenberg

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	heisenbergsimulation "posithron/x/heisenberg/simulation"
	"posithron/x/heisenberg/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	heisenbergGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&heisenbergGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCheckUncertainty          = "op_weight_msg_heisenberg"
		defaultWeightMsgCheckUncertainty int = 100
	)

	var weightMsgCheckUncertainty int
	simState.AppParams.GetOrGenerate(opWeightMsgCheckUncertainty, &weightMsgCheckUncertainty, nil,
		func(_ *rand.Rand) {
			weightMsgCheckUncertainty = defaultWeightMsgCheckUncertainty
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCheckUncertainty,
		heisenbergsimulation.SimulateMsgCheckUncertainty(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
