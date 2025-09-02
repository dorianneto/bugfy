package service

import (
	"context"
	"fmt"
	"log"
	"time"

	model "github.com/dorianneto/bugfy/internal/api/model"
	repo "github.com/dorianneto/bugfy/internal/repository"
)

type ProjectService struct {
	projectRepo *repo.ProjectRepository
	timeout     time.Duration
}

func NewProjectService(projectRepo *repo.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		timeout:     time.Duration(2) * time.Second,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, req model.RequestCreateProject) (*model.ResponseCreateProject, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("ProjectService.CreateProject - Starting project creation for: %s", req.Title)

	if req.Title == "" {
		log.Printf("ProjectService.CreateProject - Validation failed: missing required fields")
		return nil, fmt.Errorf("title is required")
	}

	p := &repo.Project{
		Title:     req.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	project, err := s.projectRepo.CreateProject(ctx, p)
	if err != nil {
		log.Printf("ProjectService.CreateProject - Database error: %v", err)
		return nil, fmt.Errorf("failed to create project: %v", err)
	}

	log.Printf("ProjectService.CreateProject - Project created successfully in database: %s", project.ID.String())

	return &model.ResponseCreateProject{
		ID:        project.ID.String(),
		Title:     project.Title,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}, nil
}

// func (s *ProjectService) DeleteUser(ctx context.Context, id uuid.UUID) error {
// 	return s.projectRepo.DeleteUser(ctx, id)
// }

// func (s *ProjectService) UpdateUsername(ctx context.Context, userID string, newUsername string) (*model.ResponseLoginUser, error) {
// 	ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()

// 	uid, err := uuid.Parse(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := s.projectRepo.UpdateUsername(ctx, uid, newUsername)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.ResponseLoginUser{
// 		ID:       user.ID.String(),
// 		Username: user.Username,
// 	}, nil
// }
