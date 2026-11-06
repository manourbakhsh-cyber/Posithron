package keeper

import (
	"context"
	"errors"
	"fmt"

	"posithron/x/entanglement/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePairRecord(ctx context.Context, msg *types.MsgCreatePairRecord) (*types.MsgCreatePairRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.PairRecord.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var pairRecord = types.PairRecord{
		Creator:    msg.Creator,
		Index:      msg.Index,
		AiNode:     msg.AiNode,
		EnergyNode: msg.EnergyNode,
		Status:     msg.Status,
	}

	if err := k.PairRecord.Set(ctx, pairRecord.Index, pairRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreatePairRecordResponse{}, nil
}

func (k msgServer) UpdatePairRecord(ctx context.Context, msg *types.MsgUpdatePairRecord) (*types.MsgUpdatePairRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.PairRecord.Get(ctx, msg.Index)
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

	var pairRecord = types.PairRecord{
		Creator:    msg.Creator,
		Index:      msg.Index,
		AiNode:     msg.AiNode,
		EnergyNode: msg.EnergyNode,
		Status:     msg.Status,
	}

	if err := k.PairRecord.Set(ctx, pairRecord.Index, pairRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update pairRecord")
	}

	return &types.MsgUpdatePairRecordResponse{}, nil
}

func (k msgServer) DeletePairRecord(ctx context.Context, msg *types.MsgDeletePairRecord) (*types.MsgDeletePairRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.PairRecord.Get(ctx, msg.Index)
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

	if err := k.PairRecord.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove pairRecord")
	}

	return &types.MsgDeletePairRecordResponse{}, nil
}
