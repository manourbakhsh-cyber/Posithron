package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"posithron/x/entanglement/keeper"
	"posithron/x/entanglement/types"
)

func createNPairRecord(keeper keeper.Keeper, ctx context.Context, n int) []types.PairRecord {
	items := make([]types.PairRecord, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].AiNode = strconv.Itoa(i)
		items[i].EnergyNode = strconv.Itoa(i)
		items[i].Status = strconv.Itoa(i)
		_ = keeper.PairRecord.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestPairRecordQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPairRecord(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPairRecordRequest
		response *types.QueryGetPairRecordResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPairRecordRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetPairRecordResponse{PairRecord: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPairRecordRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetPairRecordResponse{PairRecord: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPairRecordRequest{
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
			response, err := qs.GetPairRecord(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestPairRecordQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPairRecord(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPairRecordRequest {
		return &types.QueryAllPairRecordRequest{
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
			resp, err := qs.ListPairRecord(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PairRecord), step)
			require.Subset(t, msgs, resp.PairRecord)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListPairRecord(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PairRecord), step)
			require.Subset(t, msgs, resp.PairRecord)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListPairRecord(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.PairRecord)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListPairRecord(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
