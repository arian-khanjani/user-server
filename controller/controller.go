package controller

import (
	"context"
	pb "jonathanthegreat/mongo-repo/gen/service"
)

type repository interface {
	Get(ctx context.Context)
}

type Controller struct {
	pb.UnimplementedServerServiceServer
	repo repository
}

func New(repo repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) Get(ctx context.Context) {

}

func (c *Controller) Put(ctx context.Context) {

}
