package main

import (
	"context"
	"fmt"
	"github.com/arian-khanjani/mongo-repo/controller"
	pb "github.com/arian-khanjani/mongo-repo/gen/user"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "github.com/arian-khanjani/mongo-repo/repo/mongodb"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Server...")

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot load environment file")
	}

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
	log.Println("MongoDB connection established")

	defer func(repo *mongodb.Repo, ctx context.Context) {
		err := repo.Disconnect(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("MongoDB client disconnected")
	}(repo, ctx)

	indexes, err := repo.CreateIndexes(ctx, bson.D{
		{"name", 1},
		{"email", 1},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("indexes added:", indexes)

	ctrl := controller.New(repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening at 0.0.0.0:%s", port)

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, ctrl)

	serverCtx, serverStopCtx := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		log.Println("Closing MongoDB connection...")
		err := repo.Disconnect(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Shutting down...")
		server.GracefulStop()
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	if err := server.Serve(lis); err != nil {
		panic(err)
	}

	<-serverCtx.Done()
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
