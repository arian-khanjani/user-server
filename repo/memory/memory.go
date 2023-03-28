package memory

import (
	"context"
	"errors"
	pb "github.com/arian-khanjani/mongo-repo/gen/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type UserID string

var ErrNotFound = errors.New("not found")

type Repo struct {
	sync.RWMutex
	data map[UserID]*pb.User
}

func New() *Repo {
	return &Repo{data: map[UserID]*pb.User{}}
}

func (r *Repo) List(_ context.Context) (*pb.ListResponse, error) {
	r.RLock()
	defer r.RUnlock()

	var res pb.ListResponse

	for _, user := range r.data {
		res.Users = append(res.Users, user)
	}

	return &res, nil
}

func (r *Repo) Get(_ context.Context, in *pb.IDRequest) (*pb.User, error) {
	r.RLock()
	defer r.RUnlock()

	res, ok := r.data[UserID(in.Id)]
	if !ok {
		return nil, ErrNotFound
	}

	return res, nil
}

func (r *Repo) Update(_ context.Context, u *pb.User) (*pb.User, error) {
	r.RLock()
	defer r.RUnlock()

	r.data[UserID(u.Id)] = u

	return u, nil
}

func (r *Repo) Create(_ context.Context, u *pb.User) (*pb.User, error) {
	r.RLock()
	defer r.RUnlock()

	u.Id = primitive.NewObjectID().Hex()

	r.data[UserID(u.Id)] = u

	return u, nil
}

func (r *Repo) Delete(_ context.Context, id *pb.IDRequest) error {
	r.RLock()
	defer r.RUnlock()

	delete(r.data, UserID(id.Id))

	return nil
}
