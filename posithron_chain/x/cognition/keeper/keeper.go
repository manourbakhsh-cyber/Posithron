package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"posithron/x/cognition/types"
	entanglementkeeper "posithron/x/entanglement/keeper"
	greenenergykeeper "posithron/x/greenenergy/keeper"
)

type Keeper struct {
	storeService       corestore.KVStoreService
	cdc                codec.Codec
	addressCodec       address.Codec
	authority          []byte
	bankKeeper         types.BankKeeper
	greenEnergyKeeper  greenenergykeeper.Keeper
	entanglementKeeper entanglementkeeper.Keeper

	Schema     collections.Schema
	Params     collections.Item[types.Params]
	TaskRecord collections.Map[string, types.TaskRecord]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	greenEnergyKeeper greenenergykeeper.Keeper,
	entanglementKeeper entanglementkeeper.Keeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:       storeService,
		cdc:                cdc,
		addressCodec:       addressCodec,
		authority:          authority,
		bankKeeper:         bankKeeper,
		greenEnergyKeeper:  greenEnergyKeeper,
		entanglementKeeper: entanglementKeeper,

		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TaskRecord: collections.NewMap(sb, types.TaskRecordKey, "taskRecord", collections.StringKey, codec.CollValue[types.TaskRecord](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

func (k Keeper) GetAuthority() []byte {
	return k.authority
}