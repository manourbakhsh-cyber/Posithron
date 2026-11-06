package keeper

import (
	"context"
	"errors"

	"posithron/x/cognition/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListTaskRecord(ctx context.Context, req *types.QueryAllTaskRecordRequest) (*types.QueryAllTaskRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	taskRecords, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.TaskRecord,
		req.Pagination,
		func(_ string, value types.TaskRecord) (types.TaskRecord, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTaskRecordResponse{TaskRecord: taskRecords, Pagination: pageRes}, nil
}

func (q queryServer) GetTaskRecord(ctx context.Context, req *types.QueryGetTaskRecordRequest) (*types.QueryGetTaskRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.TaskRecord.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetTaskRecordResponse{TaskRecord: val}, nil
}
