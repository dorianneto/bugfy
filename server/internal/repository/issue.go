package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dorianneto/bugfy/internal/api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ISSUE_COLLECTION = "issues"

type Issue struct {
	ID          bson.ObjectID    `bson:"_id,omitempty"`
	ProjectID   bson.ObjectID    `bson:"project_id,omitempty"`
	Title       string           `bson:"title,omitempty"`
	Fingerprint string           `bson:"fingerprint,omitempty"`
	Count       int              `bson:"count"`
	FirstSeen   time.Time        `bson:"first_seen,omitempty"`
	LastSeen    time.Time        `bson:"last_seen,omitempty"`
	Status      model.IssueState `bson:"status"`
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

func (r *IssueRepository) FindIssuesByProject(ctx context.Context, projectId string) ([]Issue, error) {
	coll := r.db.Database("portobello").Collection(ISSUE_COLLECTION)

	d, err := bson.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %s", err)
	}

	result, err := coll.Find(ctx, bson.D{{Key: "project_id", Value: d}})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %s", err)
	}

	var i []Issue

	err = result.All(ctx, &i)
	if err != nil {
		return nil, fmt.Errorf("failed to decode issues: %s", err)
	}

	return i, nil
}
