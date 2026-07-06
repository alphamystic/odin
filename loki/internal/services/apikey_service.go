package services

import (
	"context"
	"fmt"
	dfn "github.com/alphamystic/odin/lib/definers"
)

type Apikey interface {
	CreateApikey(ctx context.Context, key dfn.Api) error
	ListApikeys(ctx context.Context, token string) ([]dfn.Api, error)
	ViewApikey(ctx context.Context, keyID string) (dfn.Api, error)
	UpdateApikey(ctx context.Context, key dfn.Api) error
	DeleteApikey(ctx context.Context, keyID string) error
}

type ApikeyService struct {
	SAC *ServerAPIConnector
}

func NewApikeyService(sac *ServerAPIConnector) *ApikeyService {
	return &ApikeyService{SAC: sac}
}

func (s *ApikeyService) CreateApikey(ctx context.Context, key dfn.Api) error {
    _, err := s.SAC.Post("/api/keys/create", key)
    if err != nil {
        return fmt.Errorf("failed to store apikey: %w", err)
    }
    return nil
}


// apikey_service.go

func (s *ApikeyService) ListApikeys(ctx context.Context, token string) ([]dfn.Api, error) {
	// Call AuthGet instead of Get
	resp, err := s.SAC.AuthGet("/api/keys/list", token)
	fmt.Println(err)
	fmt.Println(resp)
	if err != nil {
		return nil, err
	}

	if status, ok := resp["status"].(string); !ok || status != "success" {
		return nil, fmt.Errorf("API error: %v", resp["message"])
	}

	var keys []dfn.Api
	if err := s.SAC.Decode(resp["data"], &keys); err != nil {
		return nil, err
	}
	return keys, nil
}


func (s *ApikeyService) ViewApikey(ctx context.Context, keyID string) (dfn.Api, error) {
    // Fetches the specific key to populate the Reveal modal
    endpoint := fmt.Sprintf("/api/keys/view?key_id=%s", keyID)
    data, err := s.SAC.Get(endpoint)
    if err != nil {
        return dfn.Api{}, err
    }
    var key dfn.Api
    if err := s.SAC.Decode(data, &key); err != nil {
        return dfn.Api{}, err
    }
    return key, nil
}