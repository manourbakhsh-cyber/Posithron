package simulation

import (
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"posithron/x/greenenergy/keeper"
	"posithron/x/greenenergy/types"
)

func SimulateMsgCreateEnergyRecord(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Int()
		msg := &types.MsgCreateEnergyRecord{
			Creator: simAccount.Address.String(),
			Index:   strconv.Itoa(i),
		}

		found, err := k.EnergyRecord.Has(ctx, msg.Index)
		if err == nil && found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "EnergyRecord already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txGen,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUpdateEnergyRecord(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount   = simtypes.Account{}
			energyRecord = types.EnergyRecord{}
			msg          = &types.MsgUpdateEnergyRecord{}
			found        = false
		)

		var allEnergyRecord []types.EnergyRecord
		err := k.EnergyRecord.Walk(ctx, nil, func(key string, value types.EnergyRecord) (stop bool, err error) {
			allEnergyRecord = append(allEnergyRecord, value)
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		for _, obj := range allEnergyRecord {
			acc, err := ak.AddressCodec().StringToBytes(obj.Creator)
			if err != nil {
				return simtypes.OperationMsg{}, nil, err
			}

			simAccount, found = simtypes.FindAccount(accs, sdk.AccAddress(acc))
			if found {
				energyRecord = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "energyRecord creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()
		msg.Index = energyRecord.Index

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txGen,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgDeleteEnergyRecord(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount   = simtypes.Account{}
			energyRecord = types.EnergyRecord{}
			msg          = &types.MsgUpdateEnergyRecord{}
			found        = false
		)

		var allEnergyRecord []types.EnergyRecord
		err := k.EnergyRecord.Walk(ctx, nil, func(key string, value types.EnergyRecord) (stop bool, err error) {
			allEnergyRecord = append(allEnergyRecord, value)
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		for _, obj := range allEnergyRecord {
			acc, err := ak.AddressCodec().StringToBytes(obj.Creator)
			if err != nil {
				return simtypes.OperationMsg{}, nil, err
			}

			simAccount, found = simtypes.FindAccount(accs, sdk.AccAddress(acc))
			if found {
				energyRecord = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "energyRecord creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()
		msg.Index = energyRecord.Index

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txGen,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
