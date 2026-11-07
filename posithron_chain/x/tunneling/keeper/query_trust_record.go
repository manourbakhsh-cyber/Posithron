package keeper

import (
	"context"
	"errors"

	"posithron/x/tunneling/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListTrustRecord(ctx context.Context, req *types.QueryAllTrustRecordRequest) (*types.QueryAllTrustRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	trustRecords, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.TrustRecord,
		req.Pagination,
		func(_ string, value types.TrustRecord) (types.TrustRecord, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTrustRecordResponse{TrustRecord: trustRecords, Pagination: pageRes}, nil
}

func (q queryServer) GetTrustRecord(ctx context.Context, req *types.QueryGetTrustRecordRequest) (*types.QueryGetTrustRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.TrustRecord.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetTrustRecordResponse{TrustRecord: val}, nil
}
