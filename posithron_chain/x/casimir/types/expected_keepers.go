package types

import (
	"context"
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/x/simulation"
)

type AuthKeeper interface {
	simtypes.AccountKeeper
	AddressCodec() address.Codec
}

type BankKeeper interface {
	simtypes.BankKeeper
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
