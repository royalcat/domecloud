package dindex

import (
	"context"
	"dmch-server/src/domefs/media"
	"path"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoInfoIndex struct {
	coll *mongo.Collection
}

func NewVideoInfoIndex(db *mongo.Database) *VideoInfoIndex {
	coll := db.Collection("index_video_info")
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

	return &VideoInfoIndex{
		coll: coll,
	}
}

var _upsertOpts = options.Replace().SetUpsert(true)

func (vii *VideoInfoIndex) Set(v media.VisualMediaInfo) error {
	_, err := vii.coll.ReplaceOne(
		context.TODO(),
		bson.D{{Key: "path", Value: v.Path}},
		v,
		_upsertOpts,
	)
	return err
}

func (vii *VideoInfoIndex) GetSortedByDuration(targetDir string, recursive bool) ([]media.VisualMediaInfo, error) {
	ctx := context.TODO()
	targetDir = path.Clean(targetDir)
	findOpts := options.Find().SetSort(bson.D{{Key: "duration", Value: 1}})
	cur, err := vii.coll.Find(ctx,
		bson.D{{
			Key:   "path",
			Value: bson.E{Key: "$regex", Value: targetDir + "/*"},
		}}, findOpts)
	if err != nil {
		return nil, err
	}
	result := []media.VisualMediaInfo{}
	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
