package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ISSUE_COLLECTION = "issues"

type IssueState int

const (
	StateUnresolved IssueState = iota
	StateResolved
	StateIgnored
)

type Issue struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	ProjectID   bson.ObjectID `bson:"project_id,omitempty"`
	Title       string        `bson:"title,omitempty"`
	Fingerprint string        `bson:"fingerprint,omitempty"`
	Count       int           `bson:"count"`
	FirstSeen   time.Time     `bson:"first_seen,omitempty"`
	LastSeen    time.Time     `bson:"last_seen,omitempty"`
	Status      IssueState    `bson:"status"`
}

type IssueRepository struct {
	db *mongo.Client
}

func NewIssueRepository(db *mongo.Client) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) FindIssue(ctx context.Context, fingerprint string) (*Issue, error) {
	coll := r.db.Database("portobello").Collection(ISSUE_COLLECTION)

	result := coll.FindOne(ctx, bson.D{{Key: "fingerprint", Value: fingerprint}})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to find issue: %s", result.Err().Error())
	}

	var i Issue

	err := result.Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *IssueRepository) UpsertIssue(ctx context.Context, issue *Issue) (*Issue, error) {
	coll := r.db.Database("portobello").Collection(ISSUE_COLLECTION)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "fingerprint", Value: issue.Fingerprint}}
	update := bson.D{{Key: "$set", Value: issue}}

	result := coll.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to upsert issue: %s", result.Err().Error())
	}

	var i Issue

	err := result.Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}
