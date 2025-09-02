package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const PROJECT_COLLECTION = "projects"

type Project struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Title     string        `bson:"title,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}

type ProjectRepository struct {
	db *mongo.Client
}

func NewProjectRepository(db *mongo.Client) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// func (r *ProjectRepository) GetUserByID(ctx context.Context, id bson.ObjectID) (*User, error) {
// 	coll := r.db.Database("portobello").Collection("users")

// 	var user User

// 	err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
// 	if err == mongo.ErrNoDocuments {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

func (r *ProjectRepository) CreateProject(ctx context.Context, project *Project) (*Project, error) {
	coll := r.db.Database("portobello").Collection(PROJECT_COLLECTION)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "title", Value: project.Title}}
	update := bson.D{{Key: "$set", Value: project}}

	result := coll.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		return nil, fmt.Errorf("insert project: %s", result.Err().Error())
	}

	var p Project

	err := result.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// func (r *ProjectRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
// 	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
// 	if err != nil {
// 		return fmt.Errorf("delete user: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("get rows affected: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("user not found")
// 	}

// 	return nil
// }

// func (r *ProjectRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) (*User, error) {
// 	query := `
// 		UPDATE users
// 		SET username = $1, updated_at = NOW()
// 		WHERE id = $2
// 		RETURNING id, username, email, password_hash, created_at, updated_at
// 	`

// 	var user User
// 	err := r.db.QueryRowContext(ctx, query, username, id).Scan(
// 		&user.ID,
// 		&user.Username,
// 		&user.Email,
// 		&user.PasswordHash,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errors.New("user not found")
// 		}
// 		// if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
// 		// 	return nil, errors.New("username already exists")
// 		// }
// 		return nil, fmt.Errorf("update username: %w", err)
// 	}

// 	return &user, nil
// }
