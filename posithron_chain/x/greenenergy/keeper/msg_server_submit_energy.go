package keeper

import (
	"context"

	"posithron/x/greenenergy/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SubmitEnergy(ctx context.Context, msg *types.MsgSubmitEnergy) (*types.MsgSubmitEnergyResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSubmitEnergyResponse{}, nil
}
