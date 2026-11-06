package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRecordEnergy{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitEnergy{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateEnergyRecord{},
		&MsgUpdateEnergyRecord{},
		&MsgDeleteEnergyRecord{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
