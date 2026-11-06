package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		PairRecordMap: []PairRecord{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	pairRecordIndexMap := make(map[string]struct{})

	for _, elem := range gs.PairRecordMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := pairRecordIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pairRecord")
		}
		pairRecordIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
