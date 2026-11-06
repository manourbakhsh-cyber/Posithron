package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"posithron/x/tunneling/types"
)

func (k msgServer) CheckTunnel(goCtx context.Context, msg *types.MsgCheckTunnel) (*types.MsgCheckTunnelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	status := "blocked"
	if msg.TrustScore >= msg.RequiredScore {
		status = "tunneled"
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"tunnel_status_checked",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("trust_score", fmt.Sprintf("%d", msg.TrustScore)),
			sdk.NewAttribute("required_score", fmt.Sprintf("%d", msg.RequiredScore)),
			sdk.NewAttribute("status", status),
		),
	)

	return &types.MsgCheckTunnelResponse{}, nil
}
