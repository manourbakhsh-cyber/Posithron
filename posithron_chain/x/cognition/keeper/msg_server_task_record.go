package keeper

import (
	"context"
	"errors"
	"fmt"

	"posithron/x/cognition/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTaskRecord(ctx context.Context, msg *types.MsgCreateTaskRecord) (*types.MsgCreateTaskRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.TaskRecord.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var taskRecord = types.TaskRecord{
		Creator: msg.Creator,
		Index:   msg.Index,
		Prompt:  msg.Prompt,
		Bounty:  msg.Bounty,
		Status:  msg.Status,
	}

	if err := k.TaskRecord.Set(ctx, taskRecord.Index, taskRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateTaskRecordResponse{}, nil
}

func (k msgServer) UpdateTaskRecord(ctx context.Context, msg *types.MsgUpdateTaskRecord) (*types.MsgUpdateTaskRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.TaskRecord.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var taskRecord = types.TaskRecord{
		Creator: msg.Creator,
		Index:   msg.Index,
		Prompt:  msg.Prompt,
		Bounty:  msg.Bounty,
		Status:  msg.Status,
	}

	if err := k.TaskRecord.Set(ctx, taskRecord.Index, taskRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update taskRecord")
	}

	return &types.MsgUpdateTaskRecordResponse{}, nil
}

func (k msgServer) DeleteTaskRecord(ctx context.Context, msg *types.MsgDeleteTaskRecord) (*types.MsgDeleteTaskRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.TaskRecord.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.TaskRecord.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove taskRecord")
	}

	return &types.MsgDeleteTaskRecordResponse{}, nil
}
