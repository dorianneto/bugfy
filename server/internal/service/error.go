package service

import (
	"context"
	"fmt"
	"log"
	"time"

	model "github.com/dorianneto/bugfy/internal/api/model"
	repo "github.com/dorianneto/bugfy/internal/repository"
	"github.com/dorianneto/bugfy/util"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ErrorService struct {
	errorRepo *repo.ErrorRepository
	timeout   time.Duration
}

func NewErrorService(errorRepo *repo.ErrorRepository) *ErrorService {
	return &ErrorService{
		errorRepo: errorRepo,
		timeout:   time.Duration(2) * time.Second,
	}
}

func (s *ErrorService) CreateError(ctx context.Context, req model.RequestCreateError) (*model.ResponseCreateError, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("ErrorService.CreateError - Starting error creation for projectID: %s", req.ProjectID)

	pID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		log.Printf("ErrorService.CreateError - convertion error: %v", err)
		return nil, fmt.Errorf("failed to create error: %v", err)
	}

	p := &repo.Error{
		ProjectID:   pID,
		Message:     req.Message,
		Type:        "error",
		Fingerprint: util.GenerateFingerprint(req.Message),
		Context:     req.Context,
		Timestamp:   time.Now(),
	}

	e, err := s.errorRepo.CreateError(ctx, p)
	if err != nil {
		log.Printf("ErrorService.CreateError - Database error: %v", err)
		return nil, fmt.Errorf("failed to create error: %v", err)
	}

	log.Printf("ErrorService.CreateError - Project created successfully in database: %s", e.ID.String())

	return &model.ResponseCreateError{
		ID:          e.ID.String(),
		ProjectID:   e.ProjectID.String(),
		Fingerprint: e.Fingerprint,
		Message:     e.Message,
		Type:        e.Type,
		Context:     e.Context,
		Timestamp:   e.Timestamp,
	}, nil
}
