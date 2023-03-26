package main

import (
	"context"
	"jonathanthegreat/mongo-repo/controller"
	"log"
)
import "jonathanthegreat/mongo-repo/repo/mongodb"

func main() {
	repo, err := mongodb.New()
	if err != nil {
		log.Fatalln(err)
	}

	ctrl := controller.New(repo)

	ctrl.Get(context.Background())
}
