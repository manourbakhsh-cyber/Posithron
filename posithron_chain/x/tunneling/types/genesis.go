package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		TrustRecordMap: []TrustRecord{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	trustRecordIndexMap := make(map[string]struct{})

	for _, elem := range gs.TrustRecordMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := trustRecordIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for trustRecord")
		}
		trustRecordIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
