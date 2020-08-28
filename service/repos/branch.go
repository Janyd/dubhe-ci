package repos

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/service/rpc/pb"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

func NewBranchService(
	branchStore core.BranchStore,
	triggerService core.TriggerService,
) *BranchService {
	return &BranchService{
		branchStore:    branchStore,
		triggerService: triggerService,
	}
}

type BranchService struct {
	branchStore    core.BranchStore
	triggerService core.TriggerService
}

func (b *BranchService) Register(grpcServer *grpc.Server) {
	pb.RegisterBranchServiceServer(grpcServer, b)
}

func (b *BranchService) List(ctx context.Context, id *pb.Id) (*pb.Branches, error) {
	items, err := b.branchStore.List(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	branches := make([]*pb.Branch, 0)
	err = copier.Copy(&branches, items)
	if err != nil {
		return nil, errors.Error(err)
	}

	return &pb.Branches{Records: branches}, nil
}
