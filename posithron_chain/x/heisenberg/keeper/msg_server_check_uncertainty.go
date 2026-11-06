package keeper

import (
	"context"
	"fmt"

	"posithron/x/heisenberg/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CheckUncertainty(goCtx context.Context, msg *types.MsgCheckUncertainty) (*types.MsgCheckUncertaintyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ai := msg.AiShare
	energy := msg.EnergyShare

	if ai > 100 {
		ai = 100
	}
	if energy > 100 {
		energy = 100
	}

	if ai > 70 {
		energy = 100 - ai
	} else if energy > 70 {
		ai = 100 - energy
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"uncertainty_applied",
			sdk.NewAttribute("original_ai", fmt.Sprintf("%d", msg.AiShare)),
			sdk.NewAttribute("original_energy", fmt.Sprintf("%d", msg.EnergyShare)),
			sdk.NewAttribute("adjusted_ai", fmt.Sprintf("%d", ai)),
			sdk.NewAttribute("adjusted_energy", fmt.Sprintf("%d", energy)),
		),
	)

	return &types.MsgCheckUncertaintyResponse{}, nil
}
