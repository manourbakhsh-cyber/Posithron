package cognition

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"posithron/x/cognition/keeper"
	"posithron/x/cognition/types"
	entanglementkeeper "posithron/x/entanglement/keeper"
	greenenergykeeper "posithron/x/greenenergy/keeper"
	tunnelingkeeper "posithron/x/tunneling/keeper"
)

var _ depinject.OnePerModuleType = AppModule{}

func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper         types.AuthKeeper
	BankKeeper         types.BankKeeper
	GreenEnergyKeeper  greenenergykeeper.Keeper
	EntanglementKeeper entanglementkeeper.Keeper
	TunnelingKeeper    tunnelingkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	CognitionKeeper keeper.Keeper
	Module          appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.AddressCodec,
		authority,
		in.BankKeeper,
		in.GreenEnergyKeeper,
		in.EntanglementKeeper,
		in.TunnelingKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{CognitionKeeper: k, Module: m}
}