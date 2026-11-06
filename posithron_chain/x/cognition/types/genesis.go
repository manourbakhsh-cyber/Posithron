package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		TaskRecordMap: []TaskRecord{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	taskRecordIndexMap := make(map[string]struct{})

	for _, elem := range gs.TaskRecordMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := taskRecordIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for taskRecord")
		}
		taskRecordIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
