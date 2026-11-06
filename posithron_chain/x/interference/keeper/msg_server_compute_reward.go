package keeper

import (		"context"
		"fmt"
		"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"posithron/x/interference/types"
)

func (k msgServer) ComputeReward(goCtx context.Context, msg *types.MsgComputeReward) (*types.MsgComputeRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	rad := float64(msg.Phase) * math.Pi / 180.0
	a := math.Max(0, float64(msg.AiShare))
	b := math.Max(0, float64(msg.EnergyShare))

	interference := 2.0 * math.Sqrt(a*b) * math.Cos(rad)
	rewardInt := int64(interference * 1000)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"compute_reward",
			sdk.NewAttribute("task_id", msg.TaskId),
			sdk.NewAttribute("interference_reward_x1000", fmt.Sprintf("%d", rewardInt)),
		),
	)

	return &types.MsgComputeRewardResponse{}, nil
}