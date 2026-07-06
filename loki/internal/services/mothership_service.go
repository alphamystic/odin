package services

import (
	"context"
	"fmt"
	dfn "github.com/alphamystic/odin/lib/definers"
)

type MothershipSvc interface {
	CreateMS(ctx context.Context, ms dfn.Mothership, token string) error
	ListMS(ctx context.Context, active bool, token string, limit,offset int) ([]dfn.Mothership, error)
	ViewMS(ctx context.Context, msid string, token string) (dfn.Mothership, error)
	UpdateMS(ctx context.Context, ms dfn.Mothership, token string) error
}

type MothershipService struct {
	SAC *ServerAPIConnector
}

func NewMothershipService(sac *ServerAPIConnector) *MothershipService {
	return &MothershipService{SAC: sac}
}

func (s *MothershipService) ListMS(ctx context.Context, active bool, token string, limit, offset int) ([]dfn.Mothership, error) {
    // CRITICAL FIX: If your SAC returns (map, error), and you get this error,
    // it means the JSON body itself is just a number (unlikely)
    // OR your SAC signature is actually: func AuthGet(...) (int, map[string]interface{}, error)

    // Attempting to catch the status code if your SAC signature includes it:
    path := "/api/mothership/list"
    query := fmt.Sprintf("?active=%t&limit=%d&offset=%d", active, limit, offset)
    endpoint := fmt.Sprintf("%s?%s", path, query)
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return nil, err
    }

    // Double check if resp is nil or not a map
    if resp == nil {
        return nil, fmt.Errorf("empty response from API")
    }

    if status, ok := resp["status"].(string); !ok || status != "success" {
        return nil, fmt.Errorf("API error: %v", resp["message"])
    }

    var mss []dfn.Mothership
    // Decode specifically the "data" key from the response envelope
    if err := s.SAC.Decode(resp["data"], &mss); err != nil {
        fmt.Println(resp["data"])
        return nil, fmt.Errorf("failed to decode motherships: %w", err)
    }

    return mss, nil
}

func (s *MothershipService) ViewMS(ctx context.Context, msid string, token string) (dfn.Mothership, error) {
    resp, err := s.SAC.AuthGet(fmt.Sprintf("/api/mothership/view?msid=%s", msid), token)
    if err != nil {
        return dfn.Mothership{}, err
    }
    var ms dfn.Mothership
    err = s.SAC.Decode(resp["data"], &ms)
    return ms, err
}
