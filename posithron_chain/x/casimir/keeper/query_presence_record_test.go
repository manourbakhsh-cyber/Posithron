package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"posithron/x/casimir/keeper"
	"posithron/x/casimir/types"
)

func createNPresenceRecord(keeper keeper.Keeper, ctx context.Context, n int) []types.PresenceRecord {
	items := make([]types.PresenceRecord, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Status = strconv.Itoa(i)
		items[i].LastSeen = int64(i)
		_ = keeper.PresenceRecord.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestPresenceRecordQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPresenceRecord(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPresenceRecordRequest
		response *types.QueryGetPresenceRecordResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPresenceRecordRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetPresenceRecordResponse{PresenceRecord: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPresenceRecordRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetPresenceRecordResponse{PresenceRecord: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPresenceRecordRequest{
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
			response, err := qs.GetPresenceRecord(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestPresenceRecordQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPresenceRecord(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPresenceRecordRequest {
		return &types.QueryAllPresenceRecordRequest{
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
			resp, err := qs.ListPresenceRecord(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PresenceRecord), step)
			require.Subset(t, msgs, resp.PresenceRecord)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListPresenceRecord(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PresenceRecord), step)
			require.Subset(t, msgs, resp.PresenceRecord)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListPresenceRecord(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.PresenceRecord)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListPresenceRecord(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
