package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"posithron/x/hawking/types"
)

func (k msgServer) EvaporateStake(goCtx context.Context, msg *types.MsgEvaporateStake) (*types.MsgEvaporateStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nodeRecord, err := k.NodeRecord.Get(ctx, msg.NodeAddress)
	if err != nil {
		return nil, fmt.Errorf("Node not found: %s", msg.NodeAddress)
	}

	maxAllowed := msg.NetworkTotalStake / 10

	if nodeRecord.StakeAmount > maxAllowed {
		excess := nodeRecord.StakeAmount - maxAllowed
		evaporation := excess / 100

		nodeRecord.StakeAmount -= evaporation

		if err := k.NodeRecord.Set(ctx, nodeRecord.Index, nodeRecord); err != nil {
			return nil, err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"hawking_evaporation",
				sdk.NewAttribute("node_address", msg.NodeAddress),
				sdk.NewAttribute("evaporated_amount", fmt.Sprintf("%d", evaporation)),
			),
		)
	}

	return &types.MsgEvaporateStakeResponse{}, nil
}
