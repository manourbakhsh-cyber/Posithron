package keeper

import (
    "context"
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "posithron/x/casimir/types"
)

func (k msgServer) ClaimCasimir(goCtx context.Context, msg *types.MsgClaimCasimir) (*types.MsgClaimCasimirResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    record, err := k.PresenceRecord.Get(ctx, msg.NodeAddress)
    if err != nil {
        return nil, fmt.Errorf("Node not found in presence records: %s", msg.NodeAddress)
    }

    if record.Status != "active" {
        return nil, fmt.Errorf("Node is not active")
    }

    reward := sdk.NewCoins(sdk.NewInt64Coin("stake", 10))

    if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, reward); err != nil {
        return nil, err
    }

    addr, err := sdk.AccAddressFromBech32(msg.NodeAddress)
    if err != nil {
        return nil, err
    }

    if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, reward); err != nil {
        return nil, err
    }

    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "casimir_reward_claimed",
            sdk.NewAttribute("node_address", msg.NodeAddress),
            sdk.NewAttribute("reward_amount", "10stake"),
        ),
    )

    return &types.MsgClaimCasimirResponse{}, nil
}
