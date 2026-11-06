package keeper

import (
	"context"
	"errors"

	"posithron/x/casimir/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListPresenceRecord(ctx context.Context, req *types.QueryAllPresenceRecordRequest) (*types.QueryAllPresenceRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	presenceRecords, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.PresenceRecord,
		req.Pagination,
		func(_ string, value types.PresenceRecord) (types.PresenceRecord, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPresenceRecordResponse{PresenceRecord: presenceRecords, Pagination: pageRes}, nil
}

func (q queryServer) GetPresenceRecord(ctx context.Context, req *types.QueryGetPresenceRecordRequest) (*types.QueryGetPresenceRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.PresenceRecord.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetPresenceRecordResponse{PresenceRecord: val}, nil
}
