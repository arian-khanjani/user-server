package controller

import (
	"context"
	pb "github.com/jonathanthegreat/mongo-repo/gen/user"
	"go.mongodb.org/mongo-driver/bson"
)

type repository interface {
	List(context.Context) (*pb.ListResponse, error)
	Get(context.Context, *pb.IDRequest) (*pb.User, error)
	Update(context.Context, *pb.User) (*pb.User, error)
	Create(context.Context, *pb.User) (*pb.User, error)
	Delete(context.Context, *pb.IDRequest) error
	CreateIndexes(context.Context, bson.D) ([]string, error)
}

type Controller struct {
	pb.UnimplementedUserServiceServer
	repo repository
}

func New(repo repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) List(ctx context.Context, _ *pb.Empty) (*pb.ListResponse, error) {
	res, err := c.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) Get(ctx context.Context, in *pb.IDRequest) (*pb.User, error) {
	res, err := c.repo.Get(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) Update(ctx context.Context, in *pb.User) (*pb.User, error) {
	res, err := c.repo.Update(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) Create(ctx context.Context, in *pb.User) (*pb.User, error) {
	res, err := c.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) Delete(ctx context.Context, in *pb.IDRequest) (*pb.Empty, error) {
	err := c.repo.Delete(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
