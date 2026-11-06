package keeper

import (
	"context"
	"errors"
	"fmt"

	"posithron/x/casimir/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePresenceRecord(ctx context.Context, msg *types.MsgCreatePresenceRecord) (*types.MsgCreatePresenceRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.PresenceRecord.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var presenceRecord = types.PresenceRecord{
		Creator:  msg.Creator,
		Index:    msg.Index,
		Status:   msg.Status,
		LastSeen: msg.LastSeen,
	}

	if err := k.PresenceRecord.Set(ctx, presenceRecord.Index, presenceRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreatePresenceRecordResponse{}, nil
}

func (k msgServer) UpdatePresenceRecord(ctx context.Context, msg *types.MsgUpdatePresenceRecord) (*types.MsgUpdatePresenceRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.PresenceRecord.Get(ctx, msg.Index)
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

	var presenceRecord = types.PresenceRecord{
		Creator:  msg.Creator,
		Index:    msg.Index,
		Status:   msg.Status,
		LastSeen: msg.LastSeen,
	}

	if err := k.PresenceRecord.Set(ctx, presenceRecord.Index, presenceRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update presenceRecord")
	}

	return &types.MsgUpdatePresenceRecordResponse{}, nil
}

func (k msgServer) DeletePresenceRecord(ctx context.Context, msg *types.MsgDeletePresenceRecord) (*types.MsgDeletePresenceRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.PresenceRecord.Get(ctx, msg.Index)
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

	if err := k.PresenceRecord.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove presenceRecord")
	}

	return &types.MsgDeletePresenceRecordResponse{}, nil
}
