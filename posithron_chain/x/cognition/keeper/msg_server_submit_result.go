package keeper

import (
	"context"
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"posithron/x/cognition/types"
	greenenergytypes "posithron/x/greenenergy/types"
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

	// 2. Sum real energy from greenenergy module for this agent
	var totalEnergy int64 = 0
	err = k.greenEnergyKeeper.EnergyRecord.Walk(ctx, nil, func(key string, value greenenergytypes.EnergyRecord) (bool, error) {
		if value.Creator == msg.Creator {
			totalEnergy += value.Amount
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	// 3. Quantum Interference Calculation
	aiShare := 100.0
	energyShare := float64(totalEnergy)
	if energyShare == 0 {
		energyShare = 1.0 // Avoid zero to prevent math errors
	}
	
	phase := 0.0
	rad := phase * math.Pi / 180.0
	interference := 2.0 * math.Sqrt(aiShare*energyShare) * math.Cos(rad)
	rewardInt := int64(interference * 1000)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"task_completed_with_quantum_reward",
			sdk.NewAttribute("task_id", msg.TaskId),
			sdk.NewAttribute("energy_used", fmt.Sprintf("%d", totalEnergy)),
			sdk.NewAttribute("quantum_reward_x1000", fmt.Sprintf("%d", rewardInt)),
		),
	)

	// 4. Transfer Bounty
	creatorAddr, err := sdk.AccAddressFromBech32(task.Creator)
	if err != nil { return nil, fmt.Errorf("Invalid creator address: %s", err) }
	agentAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil { return nil, fmt.Errorf("Invalid agent address: %s", err) }

	if len(task.Bounty) > 0 {
		coins := sdk.NewCoins(task.Bounty...)
		if err := k.bankKeeper.SendCoins(ctx, creatorAddr, agentAddr, coins); err != nil {
			return nil, fmt.Errorf("Failed to transfer bounty: %s", err)
		}
	}

	// 5. Update task status
	task.Status = "completed"
	if err := k.TaskRecord.Set(ctx, task.Index, task); err != nil {
		return nil, err
	}

	return &types.MsgSubmitResultResponse{}, nil
}