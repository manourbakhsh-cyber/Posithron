package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:            DefaultParams(),
		PresenceRecordMap: []PresenceRecord{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	presenceRecordIndexMap := make(map[string]struct{})

	for _, elem := range gs.PresenceRecordMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := presenceRecordIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for presenceRecord")
		}
		presenceRecordIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
