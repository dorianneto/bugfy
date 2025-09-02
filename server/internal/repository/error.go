package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const ERROR_COLLECTION = "errors"

// TODO: Add stack trace to error struct
type Error struct {
	ID          bson.ObjectID     `bson:"_id,omitempty"`
	ProjectID   bson.ObjectID     `bson:"project_id,omitempty"`
	Message     string            `bson:"message,omitempty"`
	Type        string            `bson:"type,omitempty"`
	Fingerprint string            `bson:"fingerprint,omitempty"`
	Context     map[string]string `bson:"context,omitempty"`
	Timestamp   time.Time         `bson:"timestamp,omitempty"`
}

type ErrorRepository struct {
	db *mongo.Client
}

func NewErrorRepository(db *mongo.Client) *ErrorRepository {
	return &ErrorRepository{db: db}
}

func (r *ErrorRepository) CreateError(ctx context.Context, e *Error) (*Error, error) {
	coll := r.db.Database("portobello").Collection(ERROR_COLLECTION)

	result, err := coll.InsertOne(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("insert e: %s", err)
	}

	e.ID = result.InsertedID.(bson.ObjectID)

	return e, nil
}
