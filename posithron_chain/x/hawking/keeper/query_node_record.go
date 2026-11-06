package keeper

import (
	"context"
	"errors"

	"posithron/x/hawking/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListNodeRecord(ctx context.Context, req *types.QueryAllNodeRecordRequest) (*types.QueryAllNodeRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	nodeRecords, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.NodeRecord,
		req.Pagination,
		func(_ string, value types.NodeRecord) (types.NodeRecord, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNodeRecordResponse{NodeRecord: nodeRecords, Pagination: pageRes}, nil
}

func (q queryServer) GetNodeRecord(ctx context.Context, req *types.QueryGetNodeRecordRequest) (*types.QueryGetNodeRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.NodeRecord.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetNodeRecordResponse{NodeRecord: val}, nil
}
