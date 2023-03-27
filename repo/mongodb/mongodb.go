package mongodb

import (
	"context"
	"errors"
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	pb "github.com/jonathanthegreat/mongo-repo/gen/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type Repo struct {
	client mongo.Client
	db     mongo.Database
	coll   mongo.Collection
}

type ConnProps struct {
	URI  string
	DB   string
	Coll string
}

func New(props ConnProps) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(props.URI)
	c, err := mongo.NewClient(clientOptions)
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := c.Database(props.DB)

	return &Repo{
		client: *c,
		db:     *db,
		coll:   *db.Collection(props.Coll),
	}, nil
}

func NewWithRegistry(props ConnProps) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tM := reflect.TypeOf(bson.M{})
	rb := bson.NewRegistryBuilder()
	r := codecs.Register(rb).RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()
	clientOptions := options.Client().ApplyURI(props.URI).SetRegistry(r)

	c, err := mongo.NewClient(clientOptions)
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := c.Database(props.DB)

	return &Repo{
		client: *c,
		db:     *db,
		coll:   *db.Collection(props.Coll),
	}, nil
}

func (r *Repo) List(ctx context.Context) (*pb.ListResponse, error) {
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var res pb.ListResponse
	err = cur.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *Repo) Get(ctx context.Context, id *pb.IDRequest) (*pb.User, error) {
	var res pb.User
	err := r.coll.FindOne(ctx, bson.D{{"_id", id.Id}}).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *Repo) Update(ctx context.Context, u *pb.User) (*pb.User, error) {
	upd, err := r.coll.UpdateOne(ctx, bson.D{{"_id", u.Id}}, u)
	if err != nil {
		return nil, err
	}

	if upd.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return u, nil
}

func (r *Repo) Create(ctx context.Context, u *pb.User) (*pb.User, error) {
	u.Id = primitive.NewObjectID().String()
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repo) Delete(ctx context.Context, id *pb.IDRequest) error {
	one, err := r.coll.DeleteOne(ctx, bson.D{{"_id", id.Id}})
	if err != nil {
		return err
	}

	if one.DeletedCount == 0 {
		return errors.New("record was not deleted. maybe ID was incorrect")
	}

	return nil
}

func (r *Repo) CreateIndexes(ctx context.Context, index bson.D) ([]string, error) {
	res, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{{Keys: index}})
	if err != nil {
		return []string{}, err
	}

	return res, nil
}
