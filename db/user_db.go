package db

import (
	"context"

	"github.com/mukeshmahato17/hrs/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, types.User) (*types.User, error)
}

type MongoDBStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoDBStore(client *mongo.Client) *MongoDBStore {
	return &MongoDBStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoDBStore) InsertUser(ctx context.Context, user types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if objId, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = objId
	}

	return &user, nil
}

func (s *MongoDBStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	users := []*types.User{}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoDBStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user *types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}
