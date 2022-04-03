package store

import (
	"context"
	"errors"

	"github.com/256dpi/lungo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
	IsAdmin  bool   `bson:"isAdmin"`
}

type UsersStore struct {
	collection lungo.ICollection
}

func NewUsersStore(db lungo.IDatabase) *UsersStore {

	store := &UsersStore{
		collection: db.Collection("users"),
	}

	store.SetUser(context.Background(), User{Name: "admin", Password: "admin", IsAdmin: true})

	return store
}

func (s *UsersStore) SetUser(ctx context.Context, user User) error {
	_, err := s.collection.ReplaceOne(ctx,
		bson.M{"_id": user.Name}, user,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) GetUser(ctx context.Context, username string) (*User, error) {
	user := &User{}

	err := s.collection.FindOne(ctx, bson.M{"_id": username}).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
