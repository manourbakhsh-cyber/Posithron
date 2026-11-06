package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"posithron/x/qrng/types"
)

func (k msgServer) ObserveRandomness(goCtx context.Context, msg *types.MsgObserveRandomness) (*types.MsgObserveRandomnessResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	seedData := fmt.Sprintf("%s-%x-%d", msg.TargetId, ctx.BlockHeader().AppHash, ctx.BlockHeight())
	hash := sha256.Sum256([]byte(seedData))
	quantumSeed := fmt.Sprintf("%x", hash)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"qrng_observed",
			sdk.NewAttribute("target_id", msg.TargetId),
			sdk.NewAttribute("quantum_seed", quantumSeed),
			sdk.NewAttribute("observer", msg.Creator),
		),
	)

	return &types.MsgObserveRandomnessResponse{}, nil
}
