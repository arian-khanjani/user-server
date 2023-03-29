package main

import (
	"context"
	"github.com/arian-khanjani/mongo-repo/controller"
	pb "github.com/arian-khanjani/mongo-repo/gen/user"
	"github.com/arian-khanjani/mongo-repo/repo/mongodb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var client pb.UserServiceClient

var ctx = context.Background()

var tmpID string

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot load environment file")
	}

	uri := getEnv("MONGO_URI", "", true)
	db := getEnv("MONGO_DB", "user-server", false)
	coll := getEnv("MONGO_COLLECTION", "users", false)

	repo, err := mongodb.New(mongodb.ConnProps{
		URI:  uri,
		DB:   db,
		Coll: coll,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//repo := memory.New()

	ctrl := controller.New(repo)

	server := grpc.NewServer()

	pb.RegisterUserServiceServer(server, ctrl)

	lis = bufconn.Listen(bufSize)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	client = pb.NewUserServiceClient(conn)

}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreate(t *testing.T) {
	in := &pb.User{
		Name:  "John Doe",
		Email: "john.doe@gmail.com",
	}

	res, err := client.Create(ctx, in)
	if err != nil {
		t.Error(err)
	}

	tmpID = res.Id

	t.Log(res)
}

func TestList(t *testing.T) {
	res, err := client.List(ctx, &pb.Empty{})
	if err != nil {
		t.Error(err)
	}

	for i, v := range res.Users {
		t.Log(i, v)
	}
}

func TestGet(t *testing.T) {
	in := &pb.IDRequest{
		Id: tmpID,
	}

	res, err := client.Get(ctx, in)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestUpdate(t *testing.T) {
	in := &pb.User{
		Id:    tmpID,
		Name:  "John Doe - Updated",
		Email: "john.doe@yahoo.com",
	}

	res, err := client.Update(ctx, in)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestDelete(t *testing.T) {
	in := &pb.IDRequest{
		Id: tmpID,
	}

	res, err := client.Delete(ctx, in)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
