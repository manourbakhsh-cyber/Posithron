package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		EnergyRecordMap: []EnergyRecord{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	energyRecordIndexMap := make(map[string]struct{})

	for _, elem := range gs.EnergyRecordMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := energyRecordIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for energyRecord")
		}
		energyRecordIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
