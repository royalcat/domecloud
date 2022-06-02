package dindex

import (
	"context"
	"dmch-server/src/domefs/entrymodel"
	"path"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EntryIndex struct {
	coll *mongo.Collection
}

func NewEntryIndex(db *mongo.Database) *EntryIndex {
	coll := db.Collection("entry_index")
	coll.Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"path": 1,
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.M{
					"duration": 1,
				},
				Options: options.Index(),
			},
			{
				Keys: bson.M{
					"size": 1,
				},
				Options: options.Index(),
			},
			{
				Keys: bson.M{
					"modTime": 1,
				},
				Options: options.Index(),
			},
		})

	return &EntryIndex{
		coll: coll,
	}
}

var _upsertOpts = options.Replace().SetUpsert(true)

func (vii *EntryIndex) Set(ctx context.Context, v entrymodel.EntryInfo) error {
	_, err := vii.coll.ReplaceOne(
		ctx,
		bson.D{bson.E{Key: "path", Value: v.Path}},
		v,
		_upsertOpts,
	)
	return err
}

func (vii *EntryIndex) GetMediaInDir(ctx context.Context, targetDir string, recursive bool) ([]entrymodel.EntryInfo, error) {
	targetDir = path.Clean(targetDir)
	// findOpts := options.Find() //.SetSort(bson.D{{Key: "duration", Value: 1}})
	targetDir = strings.TrimRight(targetDir, "/")
	cur, err := vii.coll.Find(ctx,
		bson.D{bson.E{
			Key: "path",
			Value: bson.D{bson.E{
				Key:   "$regex",
				Value: primitive.Regex{Pattern: targetDir + "/*"},
			}},
		}},
	)
	if err != nil {
		return nil, err
	}
	result := []entrymodel.EntryInfo{}
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
