package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	casimirtypes "posithron/x/casimir/types"
	"posithron/x/cognition/types"
	entanglementtypes "posithron/x/entanglement/types"
	greenenergytypes "posithron/x/greenenergy/types"
	hawkingtypes "posithron/x/hawking/types"
	tunnelingtypes "posithron/x/tunneling/types"
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

	// 4. HEISENBERG TRADE-OFF
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

	// 5. QRNG-BASED RANDOM PHASE
	seedData := fmt.Sprintf("%s-%x-%d", msg.TaskId, ctx.BlockHeader().AppHash, ctx.BlockHeight())
	seedHash := sha256.Sum256([]byte(seedData))
	phase := float64(seedHash[0] % 91) // 0-90 degrees
	rad := phase * math.Pi / 180.0

	// 6. Interference calculation
	interference := 2.0 * math.Sqrt(aiEffective*energyEffective) * math.Cos(rad)
	rewardInt := int64(interference * 1000)

	// 7. TUNNELING BONUS: If agent's trust score >= 50, apply +10% bonus
	var trustScore int64 = 0
	var trustFound bool = false
	_ = k.tunnelingKeeper.TrustRecord.Walk(ctx, nil, func(key string, value tunnelingtypes.TrustRecord) (bool, error) {
		if value.Creator == msg.Creator {
			trustScore = value.Score
			trustFound = true
			return true, nil
		}
		return false, nil
	})

	tunnelBonusApplied := false
	if trustFound && trustScore >= 50 {
		rewardInt = int64(float64(rewardInt) * 1.10)
		tunnelBonusApplied = true
	}

	// 8. HAWKING EVAPORATION: If agent's stake > 50M, apply -5% penalty
	var agentStake int64 = 0
	var stakeFound bool = false
	_ = k.hawkingKeeper.NodeRecord.Walk(ctx, nil, func(key string, value hawkingtypes.NodeRecord) (bool, error) {
		if value.Creator == msg.Creator {
			agentStake = value.StakeAmount
			stakeFound = true
			return true, nil
		}
		return false, nil
	})

	hawkingApplied := false
	if stakeFound && agentStake > 50000000 {
		rewardInt = int64(float64(rewardInt) * 0.95) // -5% penalty
		hawkingApplied = true
	}

	// 9. CASIMIR REWARD: If agent has an active presence record, apply +5% bonus
	var casimirActive bool = false
	_ = k.casimirKeeper.PresenceRecord.Walk(ctx, nil, func(key string, value casimirtypes.PresenceRecord) (bool, error) {
		if value.Creator == msg.Creator && value.Status == "active" {
			casimirActive = true
			return true, nil
		}
		return false, nil
	})

	casimirBonusApplied := false
	if casimirActive {
		rewardInt = int64(float64(rewardInt) * 1.05) // +5% bonus
		casimirBonusApplied = true
	}

	// 10. Emit event with full quantum details
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"task_completed_with_quantum_reward",
			sdk.NewAttribute("task_id", msg.TaskId),
			sdk.NewAttribute("energy_used", fmt.Sprintf("%d", totalEnergy)),
			sdk.NewAttribute("quantum_phase_deg", fmt.Sprintf("%.0f", phase)),
			sdk.NewAttribute("ai_effective", fmt.Sprintf("%.0f", aiEffective)),
			sdk.NewAttribute("energy_effective", fmt.Sprintf("%.0f", energyEffective)),
			sdk.NewAttribute("trust_score", fmt.Sprintf("%d", trustScore)),
			sdk.NewAttribute("tunnel_bonus_applied", fmt.Sprintf("%t", tunnelBonusApplied)),
			sdk.NewAttribute("agent_stake", fmt.Sprintf("%d", agentStake)),
			sdk.NewAttribute("hawking_applied", fmt.Sprintf("%t", hawkingApplied)),
			sdk.NewAttribute("casimir_active", fmt.Sprintf("%t", casimirBonusApplied)),
			sdk.NewAttribute("quantum_reward_x1000", fmt.Sprintf("%d", rewardInt)),
		),
	)

	// 11. Check for Entanglement to split bounty
	var energyNodeAddr sdk.AccAddress
	err = k.entanglementKeeper.PairRecord.Walk(ctx, nil, func(key string, value entanglementtypes.PairRecord) (bool, error) {
		if (value.AiNode == msg.Creator || value.Creator == msg.Creator) && value.Status == "entangled" {
			addr, err := sdk.AccAddressFromBech32(value.EnergyNode)
			if err == nil && len(addr) > 0 {
				energyNodeAddr = addr
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	// 12. Transfer Bounty
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

		if energyNodeAddr != nil && !energyNodeAddr.Equals(agentAddr) {
			var energyCoins sdk.Coins
			var agentCoins sdk.Coins
			for _, coin := range coins {
				energyAmt := coin.Amount.MulRaw(30).QuoRaw(100)
				agentAmt := coin.Amount.Sub(energyAmt)
				energyCoins = energyCoins.Add(sdk.NewCoin(coin.Denom, energyAmt))
				agentCoins = agentCoins.Add(sdk.NewCoin(coin.Denom, agentAmt))
			}

			if err := k.bankKeeper.SendCoins(ctx, creatorAddr, energyNodeAddr, energyCoins); err != nil {
				return nil, fmt.Errorf("Failed to transfer energy bounty: %s", err)
			}
			if err := k.bankKeeper.SendCoins(ctx, creatorAddr, agentAddr, agentCoins); err != nil {
				return nil, fmt.Errorf("Failed to transfer agent bounty: %s", err)
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"bounty_split_entangled",
					sdk.NewAttribute("agent", msg.Creator),
					sdk.NewAttribute("energy_node", energyNodeAddr.String()),
					sdk.NewAttribute("split_ratio", "70/30"),
				),
			)
		} else {
			if err := k.bankKeeper.SendCoins(ctx, creatorAddr, agentAddr, coins); err != nil {
				return nil, fmt.Errorf("Failed to transfer bounty: %s", err)
			}
		}
	}

	// 13. Update task status
	task.Status = "completed"
	if err := k.TaskRecord.Set(ctx, task.Index, task); err != nil {
		return nil, err
	}

	return &types.MsgSubmitResultResponse{}, nil
}