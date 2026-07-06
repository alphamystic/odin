package services

import (
	"context"
	"fmt"
	dfn "github.com/alphamystic/odin/lib/definers"
)

type RMM interface {
	CreateTask(ctx context.Context, task dfn.RmmTask) error
	ViewTask(ctx context.Context, taskID string) (dfn.RmmTask, error)
	ListMyTasks(ctx context.Context) ([]dfn.RmmTask, error)
	ListExecutions(ctx context.Context, taskID, status string) ([]dfn.RmmExecution, error)
}

type RMMService struct{ SAC *ServerAPIConnector }

func NewRMMService(sac *ServerAPIConnector) *RMMService {
	return &RMMService{SAC: sac}
}


func (s *RMMService) CreateTask(ctx context.Context, task dfn.RmmTask) error {
	_, err := s.SAC.Post("/api/rmm/task/create", task)
	return err
}

func (s *RMMService) ListExecutions(ctx context.Context, taskID, status string) ([]dfn.RmmExecution, error) {
	url := fmt.Sprintf("/api/rmm/executions/list?task_id=%s&verified=%s", taskID, status)
	data, err := s.SAC.Get(url)
	var execs []dfn.RmmExecution
	err = s.SAC.Decode(data, &execs)
	return execs, err
}