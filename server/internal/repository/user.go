package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Username     string        `bson:"username,omitempty"`
	Email        string        `bson:"email,omitempty"`
	PasswordHash *string       `bson:"password_hash,omitempty"`
	CreatedAt    time.Time     `bson:"created_at,omitempty"`
	UpdatedAt    time.Time     `bson:"updated_at,omitempty"`
}

type UserRepository struct {
	db *mongo.Client
}

func NewUserRepository(db *mongo.Client) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id bson.ObjectID) (*User, error) {
	coll := r.db.Database("portobello").Collection("users")

	var user User

	err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
// 	query := `
// 		SELECT id, username, email, password_hash, created_at, updated_at
// 		FROM users
// 		WHERE email = $1
// 	`

// 	var user User
// 	err := r.db.QueryRowContext(ctx, query, email).Scan(
// 		&user.ID,
// 		&user.Username,
// 		&user.Email,
// 		&user.PasswordHash,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, nil // User not found
// 		}
// 		return nil, fmt.Errorf("query user by email: %w", err)
// 	}

// 	return &user, nil
// }

func (r *UserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	coll := r.db.Database("portobello").Collection("users")

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: user.Email}, {Key: "username", Value: user.Username}}
	update := bson.D{{Key: "$set", Value: user}}

	result := coll.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		return nil, fmt.Errorf("insert user: %w", result.Err().Error())
	}

	var u User

	err := result.Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// func (r *UserRepository) CountUsers(ctx context.Context) (int, error) {
// 	var count int
// 	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
// 	if err != nil {
// 		return 0, fmt.Errorf("count users: %w", err)
// 	}
// 	return count, nil
// }

// func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
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

// func (r *UserRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) (*User, error) {
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
