package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"posithron/x/greenenergy/keeper"
	"posithron/x/greenenergy/types"
)

func createNEnergyRecord(keeper keeper.Keeper, ctx context.Context, n int) []types.EnergyRecord {
	items := make([]types.EnergyRecord, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Amount = int64(i)
		items[i].Source = strconv.Itoa(i)
		_ = keeper.EnergyRecord.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestEnergyRecordQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEnergyRecord(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetEnergyRecordRequest
		response *types.QueryGetEnergyRecordResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetEnergyRecordRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetEnergyRecordResponse{EnergyRecord: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetEnergyRecordRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetEnergyRecordResponse{EnergyRecord: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetEnergyRecordRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetEnergyRecord(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestEnergyRecordQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEnergyRecord(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllEnergyRecordRequest {
		return &types.QueryAllEnergyRecordRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListEnergyRecord(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.EnergyRecord), step)
			require.Subset(t, msgs, resp.EnergyRecord)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListEnergyRecord(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.EnergyRecord), step)
			require.Subset(t, msgs, resp.EnergyRecord)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListEnergyRecord(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.EnergyRecord)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListEnergyRecord(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
