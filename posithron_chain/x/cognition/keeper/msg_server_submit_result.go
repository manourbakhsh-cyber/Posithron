package keeper

import (
	"context"
	"crypto/sha256"
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

	// 3. Base shares
	aiShare := 100.0
	energyShare := float64(totalEnergy)
	if energyShare == 0 {
		energyShare = 1.0
	}

	// 4. HEISENBERG TRADE-OFF: Can't be best at both
	aiEffective := aiShare
	energyEffective := energyShare
	if aiShare > 70.0 && energyShare > 70.0 {
		aiEffective = aiShare * 0.85
		energyEffective = energyShare * 0.85
	} else if aiShare > 70.0 {
		energyEffective = energyShare * 0.7
	} else if energyShare > 70.0 {
		aiEffective = aiShare * 0.7
	}

	// 5. QRNG-BASED RANDOM PHASE (deterministic per block, unpredictable)
	seedData := fmt.Sprintf("%s-%x-%d", msg.TaskId, ctx.BlockHeader().AppHash, ctx.BlockHeight())
	seedHash := sha256.Sum256([]byte(seedData))
	phase := float64(seedHash[0] % 91) // 0-90 degrees (cos always >= 0)
	rad := phase * math.Pi / 180.0

	// 6. Interference calculation
	interference := 2.0 * math.Sqrt(aiEffective*energyEffective) * math.Cos(rad)
	rewardInt := int64(interference * 1000)

	// 7. Emit event with full quantum details
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"task_completed_with_quantum_reward",
			sdk.NewAttribute("task_id", msg.TaskId),
			sdk.NewAttribute("energy_used", fmt.Sprintf("%d", totalEnergy)),
			sdk.NewAttribute("quantum_phase_deg", fmt.Sprintf("%.0f", phase)),
			sdk.NewAttribute("ai_effective", fmt.Sprintf("%.0f", aiEffective)),
			sdk.NewAttribute("energy_effective", fmt.Sprintf("%.0f", energyEffective)),
			sdk.NewAttribute("quantum_reward_x1000", fmt.Sprintf("%d", rewardInt)),
		),
	)

	// 8. Transfer bounty
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

	// 9. Update task status
	task.Status = "completed"
	if err := k.TaskRecord.Set(ctx, task.Index, task); err != nil {
		return nil, err
	}

	return &types.MsgSubmitResultResponse{}, nil
}