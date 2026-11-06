package tunneling

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

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
	}
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

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
