package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"posithron/x/cognition/types"
)

func (k msgServer) CreateTask(goCtx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate a unique ID for the task using block height and creator address
	index := strconv.FormatInt(ctx.BlockHeight(), 10) + "-" + msg.Creator

	// Create the Task Record
	taskRecord := types.TaskRecord{
		Index:   index,
		Creator: msg.Creator,
		Prompt:  msg.Prompt,
		Bounty:  msg.Bounty,
		Status:  "pending", // Set initial status to pending
	}

	// Save the task into the blockchain storage using the new collections API
	if err := k.TaskRecord.Set(ctx, taskRecord.Index, taskRecord); err != nil {
		return nil, err
	}

	return &types.MsgCreateTaskResponse{}, nil
}
