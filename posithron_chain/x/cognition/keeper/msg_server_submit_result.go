package keeper

import (
	"context"
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"posithron/x/cognition/types"
)

func (k msgServer) SubmitResult(goCtx context.Context, msg *types.MsgSubmitResult) (*types.MsgSubmitResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Get the task
	task, err := k.TaskRecord.Get(ctx, msg.TaskId)
	if err != nil {
		return nil, fmt.Errorf("Task not found: %s", msg.TaskId)
	}

	if task.Status == "completed" {
		return nil, fmt.Errorf("Task already completed: %s", msg.TaskId)
	}

	// 2. Quantum Interference Calculation (MVP Test Values)
	aiShare := 100.0
	energyShare := 100.0
	phase := 0.0
	rad := phase * math.Pi / 180.0
	interference := 2.0 * math.Sqrt(aiShare*energyShare) * math.Cos(rad)
	rewardInt := int64(interference * 1000)

	// Emit the quantum reward event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"task_completed_with_quantum_reward",
			sdk.NewAttribute("task_id", msg.TaskId),
			sdk.NewAttribute("quantum_reward_x1000", fmt.Sprintf("%d", rewardInt)),
		),
	)

	// 3. Transfer the Bounty (from Creator to Agent)
	creatorAddr, err := sdk.AccAddressFromBech32(task.Creator)
	if err != nil {
		return nil, fmt.Errorf("Invalid creator address: %s", err)
	}

	agentAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, fmt.Errorf("Invalid agent address: %s", err)
	}

	if len(task.Bounty) > 0 {
		coins := sdk.NewCoins(task.Bounty...)
		if err := k.bankKeeper.SendCoins(ctx, creatorAddr, agentAddr, coins); err != nil {
			return nil, fmt.Errorf("Failed to transfer bounty: %s", err)
		}
	}

	// 4. Update the task status
	task.Status = "completed"
	if err := k.TaskRecord.Set(ctx, task.Index, task); err != nil {
		return nil, err
	}

	return &types.MsgSubmitResultResponse{}, nil
}