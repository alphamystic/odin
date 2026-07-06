package services

import (
	"context"
	"fmt"
	dfn "github.com/alphamystic/odin/lib/definers"
)

type MinionSvc interface {
	CreateMinion(ctx context.Context, minion dfn.Minion) error
	ViewMinion(ctx context.Context, minionID string) (dfn.Minion, error)
	ListMinions(ctx context.Context, limit, offset int) ([]dfn.Minion, error)
	DeactivateMinion(ctx context.Context, minionID string) error
}

type MinionService struct{ SAC *ServerAPIConnector }

func NewMinionService(sac *ServerAPIConnector) *MinionService {
	return &MinionService{SAC: sac}
}

func (s *MinionService) ListMinions(ctx context.Context, limit, offset int) ([]dfn.Minion, error) {
	data, err := s.SAC.Get(fmt.Sprintf("/api/minion/list?limit=%d&offset=%d", limit, offset))
	var mins []dfn.Minion
	err = s.SAC.Decode(data, &mins)
	return mins, err
}