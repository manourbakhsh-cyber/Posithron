package greenenergy

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/testutil/sample"
	greenenergysimulation "posithron/x/greenenergy/simulation"
	"posithron/x/greenenergy/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	greenenergyGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		EnergyRecordMap: []types.EnergyRecord{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&greenenergyGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateEnergyRecord          = "op_weight_msg_greenenergy"
		defaultWeightMsgCreateEnergyRecord int = 100
	)

	var weightMsgCreateEnergyRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEnergyRecord, &weightMsgCreateEnergyRecord, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEnergyRecord = defaultWeightMsgCreateEnergyRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEnergyRecord,
		greenenergysimulation.SimulateMsgCreateEnergyRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateEnergyRecord          = "op_weight_msg_greenenergy"
		defaultWeightMsgUpdateEnergyRecord int = 100
	)

	var weightMsgUpdateEnergyRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEnergyRecord, &weightMsgUpdateEnergyRecord, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEnergyRecord = defaultWeightMsgUpdateEnergyRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEnergyRecord,
		greenenergysimulation.SimulateMsgUpdateEnergyRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteEnergyRecord          = "op_weight_msg_greenenergy"
		defaultWeightMsgDeleteEnergyRecord int = 100
	)

	var weightMsgDeleteEnergyRecord int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteEnergyRecord, &weightMsgDeleteEnergyRecord, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteEnergyRecord = defaultWeightMsgDeleteEnergyRecord
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteEnergyRecord,
		greenenergysimulation.SimulateMsgDeleteEnergyRecord(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgSubmitEnergy          = "op_weight_msg_greenenergy"
		defaultWeightMsgSubmitEnergy int = 100
	)

	var weightMsgSubmitEnergy int
	simState.AppParams.GetOrGenerate(opWeightMsgSubmitEnergy, &weightMsgSubmitEnergy, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitEnergy = defaultWeightMsgSubmitEnergy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitEnergy,
		greenenergysimulation.SimulateMsgSubmitEnergy(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRecordEnergy          = "op_weight_msg_greenenergy"
		defaultWeightMsgRecordEnergy int = 100
	)

	var weightMsgRecordEnergy int
	simState.AppParams.GetOrGenerate(opWeightMsgRecordEnergy, &weightMsgRecordEnergy, nil,
		func(_ *rand.Rand) {
			weightMsgRecordEnergy = defaultWeightMsgRecordEnergy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRecordEnergy,
		greenenergysimulation.SimulateMsgRecordEnergy(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
