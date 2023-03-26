package mongodb

import (
	"context"
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"time"
)

type Repo struct {
	client mongo.Client
	db     mongo.Database
	coll   mongo.Collection
}

func New() (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tM := reflect.TypeOf(bson.M{})
	rb := bson.NewRegistryBuilder()
	r := codecs.Register(rb).RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()
	clientOptions := options.Client().ApplyURI("").SetRegistry(r)

	c, err := mongo.NewClient(clientOptions)
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	db := c.Database("")

	return &Repo{
		client: *c,
		db:     *db,
		coll:   *db.Collection(""),
	}, nil
}

func (r *Repo) Get(ctx context.Context) {

}

func CreateIndex(c *mongo.Collection, index bson.D) {
	model := []mongo.IndexModel{{Keys: index}}

	_, err := c.Indexes().CreateMany(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
	}
}
