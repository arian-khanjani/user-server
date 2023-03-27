package main

import (
	"fmt"
	"github.com/jonathanthegreat/mongo-repo/controller"
	pb "github.com/jonathanthegreat/mongo-repo/gen/user"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)
import "github.com/jonathanthegreat/mongo-repo/repo/mongodb"

func main() {
	uri := getEnv("MONGO_URI", "", true)
	db := getEnv("MONGO_DB", "user-server", false)
	coll := getEnv("MONGO_COLLECTION", "users", false)
	port := getEnv("APP_PORT", "3000", false)

	repo, err := mongodb.New(mongodb.ConnProps{
		URI:  uri,
		DB:   db,
		Coll: coll,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctrl := controller.New(repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, ctrl)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

func getEnv(key, def string, must bool) string {
	res := os.Getenv(key)
	if res == "" {
		if must {
			panic(fmt.Sprintf("env \"%s\" not found", key))
		}
		return def
	}

	return res
}
