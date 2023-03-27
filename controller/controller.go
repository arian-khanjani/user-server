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

func (c *Controller) Get(ctx context.Context) {

}

func (c *Controller) Put(ctx context.Context) {

}
