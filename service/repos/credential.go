package repos

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/service/rpc/pb"
	"dubhe-ci/utils"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

func NewCredentialService(credStore core.CredentialStore) *CredentialService {
	return &CredentialService{credStore: credStore}
}

type CredentialService struct {
	credStore core.CredentialStore
}

func (c *CredentialService) Register(grpcServer *grpc.Server) {
	pb.RegisterCredentialServiceServer(grpcServer, c)
}

func (c *CredentialService) List(ctx context.Context, _ *pb.Empty) (*pb.Creds, error) {
	credentials, err := c.credStore.List(ctx)
	if err != nil {
		return nil, err
	}

	creds := make([]*pb.Cred, 0)
	err = copier.Copy(&creds, credentials)
	if err != nil {
		return nil, errors.New(999999)
	}

	return &pb.Creds{Records: creds}, nil
}

func (c *CredentialService) Create(ctx context.Context, cred *pb.Cred) (*pb.Empty, error) {
	credential := &core.Credential{}
	err := copier.Copy(credential, cred)
	if err != nil {
		return nil, errors.Error(err)
	}

	_, err = c.credStore.Create(ctx, credential)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (c *CredentialService) Delete(ctx context.Context, id *pb.Id) (*pb.Empty, error) {
	credential, err := c.credStore.Find(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	err = c.credStore.Delete(ctx, credential)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (c *CredentialService) RandomGenerateSshKey(ctx context.Context, empty *pb.Empty) (*pb.Key, error) {
	privateKey, publicKey, err := utils.GenRsaKey(1024)
	if err != nil {
		return nil, err
	}

	return &pb.Key{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}
