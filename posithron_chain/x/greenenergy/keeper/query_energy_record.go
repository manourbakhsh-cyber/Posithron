package keeper

import (
	"context"
	"errors"

	"posithron/x/greenenergy/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListEnergyRecord(ctx context.Context, req *types.QueryAllEnergyRecordRequest) (*types.QueryAllEnergyRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	energyRecords, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.EnergyRecord,
		req.Pagination,
		func(_ string, value types.EnergyRecord) (types.EnergyRecord, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEnergyRecordResponse{EnergyRecord: energyRecords, Pagination: pageRes}, nil
}

func (q queryServer) GetEnergyRecord(ctx context.Context, req *types.QueryGetEnergyRecordRequest) (*types.QueryGetEnergyRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.EnergyRecord.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetEnergyRecordResponse{EnergyRecord: val}, nil
}
