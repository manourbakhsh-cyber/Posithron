package keeper

import (
	"context"
	"errors"
	"fmt"

	"posithron/x/hawking/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateNodeRecord(ctx context.Context, msg *types.MsgCreateNodeRecord) (*types.MsgCreateNodeRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.NodeRecord.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var nodeRecord = types.NodeRecord{
		Creator:     msg.Creator,
		Index:       msg.Index,
		StakeAmount: msg.StakeAmount,
	}

	if err := k.NodeRecord.Set(ctx, nodeRecord.Index, nodeRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateNodeRecordResponse{}, nil
}

func (k msgServer) UpdateNodeRecord(ctx context.Context, msg *types.MsgUpdateNodeRecord) (*types.MsgUpdateNodeRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.NodeRecord.Get(ctx, msg.Index)
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

	var nodeRecord = types.NodeRecord{
		Creator:     msg.Creator,
		Index:       msg.Index,
		StakeAmount: msg.StakeAmount,
	}

	if err := k.NodeRecord.Set(ctx, nodeRecord.Index, nodeRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update nodeRecord")
	}

	return &types.MsgUpdateNodeRecordResponse{}, nil
}

func (k msgServer) DeleteNodeRecord(ctx context.Context, msg *types.MsgDeleteNodeRecord) (*types.MsgDeleteNodeRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.NodeRecord.Get(ctx, msg.Index)
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

	if err := k.NodeRecord.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove nodeRecord")
	}

	return &types.MsgDeleteNodeRecordResponse{}, nil
}
