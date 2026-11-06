package keeper

import (
	"context"
	"errors"

	"posithron/x/entanglement/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListPairRecord(ctx context.Context, req *types.QueryAllPairRecordRequest) (*types.QueryAllPairRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pairRecords, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.PairRecord,
		req.Pagination,
		func(_ string, value types.PairRecord) (types.PairRecord, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPairRecordResponse{PairRecord: pairRecords, Pagination: pageRes}, nil
}

func (q queryServer) GetPairRecord(ctx context.Context, req *types.QueryGetPairRecordRequest) (*types.QueryGetPairRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.PairRecord.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetPairRecordResponse{PairRecord: val}, nil
}
