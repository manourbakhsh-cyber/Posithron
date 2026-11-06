package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		NodeRecordMap: []NodeRecord{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	nodeRecordIndexMap := make(map[string]struct{})

	for _, elem := range gs.NodeRecordMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := nodeRecordIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for nodeRecord")
		}
		nodeRecordIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
