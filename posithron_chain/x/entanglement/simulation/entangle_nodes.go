package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"posithron/x/entanglement/keeper"
	"posithron/x/entanglement/types"
)

func SimulateMsgEntangleNodes(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgEntangleNodes{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the EntangleNodes simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "EntangleNodes simulation not implemented"), nil, nil
	}
}
