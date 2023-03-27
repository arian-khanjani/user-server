package main

import (
	"context"
	"github.com/jonathanthegreat/mongo-repo/controller"
	"log"
)
import "github.com/jonathanthegreat/mongo-repo/repo/mongodb"

func main() {
	repo, err := mongodb.New(mongodb.ConnProps{})
	if err != nil {
		log.Fatalln(err)
	}

	ctrl := controller.New(repo)

	ctrl.Get(context.Background())
}
