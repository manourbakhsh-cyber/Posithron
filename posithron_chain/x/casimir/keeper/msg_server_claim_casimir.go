package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"posithron/x/casimir/types"
)

func (k msgServer) ClaimCasimir(goCtx context.Context, msg *types.MsgClaimCasimir) (*types.MsgClaimCasimirResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Find presence record by index (NodeAddress)
	record, err := k.PresenceRecord.Get(ctx, msg.NodeAddress)
	if err != nil {
		return nil, fmt.Errorf("Node not found in presence records: %s", msg.NodeAddress)
	}

	if record.Status != "active" {
		return nil, fmt.Errorf("Node is not active")
	}

	// 2. Define reward
	reward := sdk.NewCoins(sdk.NewInt64Coin("stake", 10))

	// 3. Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, reward); err != nil {
		return nil, err
	}

	// 4. Send coins to the creator of the presence record (the actual Bech32 address)
	addr, err := sdk.AccAddressFromBech32(record.Creator)
	if err != nil {
		return nil, err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, reward); err != nil {
		return nil, err
	}

	// 5. Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"casimir_reward_claimed",
			sdk.NewAttribute("node_address", record.Creator),
			sdk.NewAttribute("reward_amount", "10stake"),
		),
	)

	return &types.MsgClaimCasimirResponse{}, nil
}
