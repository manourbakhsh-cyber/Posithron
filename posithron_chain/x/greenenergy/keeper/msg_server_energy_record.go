package keeper

import (
	"context"
	"errors"
	"fmt"

	"posithron/x/greenenergy/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEnergyRecord(ctx context.Context, msg *types.MsgCreateEnergyRecord) (*types.MsgCreateEnergyRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.EnergyRecord.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var energyRecord = types.EnergyRecord{
		Creator: msg.Creator,
		Index:   msg.Index,
		Amount:  msg.Amount,
		Source:  msg.Source,
	}

	if err := k.EnergyRecord.Set(ctx, energyRecord.Index, energyRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateEnergyRecordResponse{}, nil
}

func (k msgServer) UpdateEnergyRecord(ctx context.Context, msg *types.MsgUpdateEnergyRecord) (*types.MsgUpdateEnergyRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.EnergyRecord.Get(ctx, msg.Index)
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

	var energyRecord = types.EnergyRecord{
		Creator: msg.Creator,
		Index:   msg.Index,
		Amount:  msg.Amount,
		Source:  msg.Source,
	}

	if err := k.EnergyRecord.Set(ctx, energyRecord.Index, energyRecord); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update energyRecord")
	}

	return &types.MsgUpdateEnergyRecordResponse{}, nil
}

func (k msgServer) DeleteEnergyRecord(ctx context.Context, msg *types.MsgDeleteEnergyRecord) (*types.MsgDeleteEnergyRecordResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.EnergyRecord.Get(ctx, msg.Index)
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

	if err := k.EnergyRecord.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove energyRecord")
	}

	return &types.MsgDeleteEnergyRecordResponse{}, nil
}
