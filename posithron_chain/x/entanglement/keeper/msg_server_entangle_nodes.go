package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"posithron/x/entanglement/types"
)

func (k msgServer) EntangleNodes(goCtx context.Context, msg *types.MsgEntangleNodes) (*types.MsgEntangleNodesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	index := strconv.FormatInt(ctx.BlockHeight(), 10) + "-" + msg.AiNode + "-" + msg.EnergyNode

	pairRecord := types.PairRecord{
		Index:      index,
		AiNode:     msg.AiNode,
		EnergyNode: msg.EnergyNode,
		Status:     "entangled",
	}

	if err := k.PairRecord.Set(ctx, pairRecord.Index, pairRecord); err != nil {
		return nil, err
	}

	return &types.MsgEntangleNodesResponse{}, nil
}
