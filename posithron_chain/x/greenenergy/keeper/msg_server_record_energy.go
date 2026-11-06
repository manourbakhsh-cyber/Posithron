package keeper

import (		"context"
		"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"posithron/x/greenenergy/types"
)

func (k msgServer) RecordEnergy(goCtx context.Context, msg *types.MsgRecordEnergy) (*types.MsgRecordEnergyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	index := strconv.FormatInt(ctx.BlockHeight(), 10) + "-" + msg.Creator

	energyRecord := types.EnergyRecord{
		Index:   index,
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Source:  msg.Source,
	}

	if err := k.EnergyRecord.Set(ctx, energyRecord.Index, energyRecord); err != nil {
		return nil, err
	}

	return &types.MsgRecordEnergyResponse{}, nil
}
