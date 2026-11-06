package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"posithron/x/cognition/keeper"
	"posithron/x/cognition/types"
)

func createNTaskRecord(keeper keeper.Keeper, ctx context.Context, n int) []types.TaskRecord {
	items := make([]types.TaskRecord, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Prompt = strconv.Itoa(i)
		items[i].Bounty = sdk.NewCoins(sdk.NewInt64Coin(`token`, int64(i%1+100)), sdk.NewInt64Coin(`stake`, int64(i%2+100)))
		items[i].Status = strconv.Itoa(i)
		_ = keeper.TaskRecord.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestTaskRecordQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNTaskRecord(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetTaskRecordRequest
		response *types.QueryGetTaskRecordResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetTaskRecordRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetTaskRecordResponse{TaskRecord: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetTaskRecordRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetTaskRecordResponse{TaskRecord: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetTaskRecordRequest{
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
			response, err := qs.GetTaskRecord(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestTaskRecordQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNTaskRecord(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTaskRecordRequest {
		return &types.QueryAllTaskRecordRequest{
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
			resp, err := qs.ListTaskRecord(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.TaskRecord), step)
			require.Subset(t, msgs, resp.TaskRecord)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListTaskRecord(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.TaskRecord), step)
			require.Subset(t, msgs, resp.TaskRecord)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListTaskRecord(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.TaskRecord)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListTaskRecord(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
