package service

import (
	"context"
	"fmt"
	"log"
	"time"

	repo "github.com/dorianneto/bugfy/internal/repository"
)

type IssueService struct {
	issueRepo *repo.IssueRepository
	timeout   time.Duration
}

func NewIssueService(issueRepo *repo.IssueRepository) *IssueService {
	return &IssueService{
		issueRepo: issueRepo,
		timeout:   time.Duration(2) * time.Second,
	}
}

func (s *IssueService) GroupError(ctx context.Context, e *repo.Error) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	log.Printf("IssueService.GroupError - Starting group error creation for fingerprint: %s", e.Fingerprint)

	var issue *repo.Issue

	issue = &repo.Issue{
		ProjectID:   e.ProjectID,
		Fingerprint: e.Fingerprint,
		Title:       e.Message,
		Count:       1,
		FirstSeen:   e.Timestamp,
		LastSeen:    e.Timestamp,
		Status:      repo.StateUnresolved,
	}

	i, err := s.issueRepo.FindIssue(ctx, e.Fingerprint)
	if err != nil {
		log.Printf("IssueService.GroupError - Database error: %v", err)
		return fmt.Errorf("failed to find issue: %v", err)
	}

	if i != nil {
		issue = i

		issue.Count++
		issue.LastSeen = e.Timestamp
	}

	_, err = s.issueRepo.UpsertIssue(ctx, issue)
	if err != nil {
		log.Printf("IssueService.GroupError - Database error: %v", err)
		return fmt.Errorf("failed to upsert issue: %v", err)
	}

	return nil
}
