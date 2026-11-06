package keeper

import (		"context"
		"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"posithron/x/cognition/types"
)

func (k msgServer) SubmitResult(goCtx context.Context, msg *types.MsgSubmitResult) (*types.MsgSubmitResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	task, err := k.TaskRecord.Get(ctx, msg.TaskId)
	if err != nil {
		return nil, fmt.Errorf("Task not found: %s", msg.TaskId)
	}

	if task.Status == "completed" {
		return nil, fmt.Errorf("Task already completed: %s", msg.TaskId)
	}

	// Parse addresses for bank transfer
	creatorAddr, err := sdk.AccAddressFromBech32(task.Creator)
	if err != nil {
		return nil, fmt.Errorf("Invalid creator address: %s", err)
	}

	agentAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, fmt.Errorf("Invalid agent address: %s", err)
	}

	// Transfer the bounty from the task creator to the agent
	if len(task.Bounty) > 0 {
		coins := sdk.NewCoins(task.Bounty...)
		if err := k.bankKeeper.SendCoins(ctx, creatorAddr, agentAddr, coins); err != nil {
			return nil, fmt.Errorf("Failed to transfer bounty: %s", err)
		}
	}

	// Update status to completed
	task.Status = "completed"
	if err := k.TaskRecord.Set(ctx, task.Index, task); err != nil {
		return nil, err
	}

	return &types.MsgSubmitResultResponse{}, nil
}
