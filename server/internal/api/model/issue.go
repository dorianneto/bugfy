package model

import (
	"time"
)

type IssueState int

const (
	IssueStateUnresolved IssueState = iota
	IssueStateResolved
	IssueStateIgnored
)

type ResponseGetIssues struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Fingerprint string     `json:"fingerprint"`
	Count       int        `json:"count"`
	FirstSeen   time.Time  `json:"first_seen"`
	LastSeen    time.Time  `json:"last_seen"`
	Status      IssueState `json:"status"`
}
