package repos

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/core"
	"dubhe-ci/service/rpc/pb"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

func NewBuildService(
	buildStore core.BuildStore,
	branchStore core.BranchStore,
	stepStore core.StepStore,
) *BuildService {
	return &BuildService{
		buildStore:  buildStore,
		branchStore: branchStore,
		stepStore:   stepStore,
	}
}

type BuildService struct {
	buildStore  core.BuildStore
	branchStore core.BranchStore
	stepStore   core.StepStore
}

func (b *BuildService) Register(grpcServer *grpc.Server) {
	pb.RegisterBuildServiceServer(grpcServer, b)
}

func (b *BuildService) List(ctx context.Context, request *pb.BuildListRequest) (*pb.BuildRecords, error) {
	page := request.Page
	page.Desc = append(page.Desc, "created")
	p := &common.Page{
		Current: page.Current,
		Size:    page.Size,
		Desc:    page.Desc,
		Asc:     page.Asc,
	}

	branch, err := b.branchStore.Find(ctx, request.BranchId)
	if err != nil {
		return nil, err
	}

	p, err = b.buildStore.List(ctx, request.RepoId, branch.Name, p)
	if err != nil {
		return nil, err
	}
	page.Total = p.Total

	builds := make([]*pb.Build, 0)
	err = copier.Copy(&builds, p.Records)
	if err != nil {
		return nil, err
	}

	return &pb.BuildRecords{
		Page:   page,
		Builds: builds,
	}, nil

}

func (b *BuildService) Find(ctx context.Context, id *pb.Id) (*pb.Build, error) {
	build, err := b.buildStore.Find(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	bd := &pb.Build{}
	err = copier.Copy(bd, build)
	if err != nil {
		return nil, err
	}

	return bd, nil
}

func (b *BuildService) Delete(ctx context.Context, id *pb.Id) (*pb.Empty, error) {
	build, err := b.buildStore.Find(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	err = b.buildStore.Delete(ctx, build)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (b *BuildService) ListStep(ctx context.Context, id *pb.Id) (*pb.Steps, error) {
	items, err := b.stepStore.List(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	steps := make([]*pb.Step, 0)

	err = copier.Copy(&steps, items)
	if err != nil {
		return nil, err
	}

	return &pb.Steps{Records: steps}, nil
}
