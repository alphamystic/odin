package services

import (
	"context"
	"fmt"
	"net/url"

	png_hnd "github.com/alphamystic/odin/lib/handlers"
)

type Recon interface {
	CreateScan(ctx context.Context, scan png_hnd.Scans) (string, error)
	ListScans(ctx context.Context, token, scanType string) ([]png_hnd.Scans, error)
	ListTargets(ctx context.Context, token, scanID string) ([]png_hnd.Target, error)
	CreateTarget(ctx context.Context, target png_hnd.Target) (string, error)
	ViewTarget(ctx context.Context, token, scanID, targetID string) (map[string]interface{}, error)
	CreateVulnerability(ctx context.Context, vuln png_hnd.Vulnerabilities) error
	ListVulnerabilities(ctx context.Context, token, targetID string) ([]png_hnd.Vulnerabilities, error)
}

type ReconService struct{ SAC *ServerAPIConnector }

func NewReconService(sac *ServerAPIConnector) *ReconService {
	return &ReconService{SAC: sac}
}

// POST {{url_endpoint}}/api/recon/createscan/
func (s *ReconService) CreateScan(ctx context.Context, scan png_hnd.Scans) (string, error) {
	resp, err := s.SAC.Post("/api/recon/createscan/", scan)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resp["redirecturl"]), nil
}

// GET {{url_endpoint}}/api/recon/listscans?scantype=ALL
func (s *ReconService) ListScans(ctx context.Context, token, scanType string) ([]png_hnd.Scans, error) {
	endpoint := fmt.Sprintf("/api/recon/listscans?scantype=%s", url.QueryEscape(scanType))
	resp, err := s.SAC.AuthGet(endpoint, token)
	if err != nil {
		return nil, err
	}
	dataField, ok := resp["data"]
	if !ok {
		return nil, fmt.Errorf("API response(ListScans) missing 'data' field")
	}
	var scans []png_hnd.Scans
	if err := s.SAC.Decode(dataField, &scans); err != nil {
		return nil, fmt.Errorf("failed to decode scans: %w", err)
	}
	return scans, nil
}

// GET {{url_endpoint}}/api/recon/viewscans?scanid=UUID
func (s *ReconService) ListTargets(ctx context.Context, token, scanID string) ([]png_hnd.Target, error) {
	endpoint := fmt.Sprintf("/api/recon/viewscans?scanid=%s", url.QueryEscape(scanID))
	resp, err := s.SAC.AuthGet(endpoint, token)
	if err != nil {
		return nil, err
	}
	dataField, ok := resp["data"]
	if !ok {
		return nil, fmt.Errorf("API response(ListTargets) missing 'data' field")
	}
	var targets []png_hnd.Target
	if err := s.SAC.Decode(dataField, &targets); err != nil {
		return nil, fmt.Errorf("failed to decode targets: %w", err)
	}
	return targets, nil
}

// POST {{url_endpoint}}/api/recon/createtarget
func (s *ReconService) CreateTarget(ctx context.Context, target png_hnd.Target) (string, error) {
	resp, err := s.SAC.Post("/api/recon/createtarget", target)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resp["redirecturl"]), nil
}

// GET {{url_endpoint}}/api/recon/viewtarget?scan-id=UUID&target-id=UUID
func (s *ReconService) ViewTarget(ctx context.Context, token, scanID, targetID string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/api/recon/viewtarget?scan-id=%s&target-id=%s", url.QueryEscape(scanID), url.QueryEscape(targetID))
	resp, err := s.SAC.AuthGet(endpoint, token)
	if err != nil {
		return nil, err
	}
	dataField, ok := resp["data"]
	if !ok {
		return nil, fmt.Errorf("API response(ViewTarget) missing 'data' field")
	}
	var targetData map[string]interface{}
	if err := s.SAC.Decode(dataField, &targetData); err != nil {
		return nil, fmt.Errorf("failed to decode target data payload: %w", err)
	}
    //fmt.Println(targetData)
	return targetData, nil
}

// POST {{url_endpoint}}/api/recon/vulnerability/create
func (s *ReconService) CreateVulnerability(ctx context.Context, vuln png_hnd.Vulnerabilities) error {
	_, err := s.SAC.Post("/api/recon/vulnerability/create", vuln)
	return err
}

// GET {{url_endpoint}}/api/recon/vulnerability/list?target_id=UUID
func (s *ReconService) ListVulnerabilities(ctx context.Context, token, targetID string) ([]png_hnd.Vulnerabilities, error) {
	endpoint := fmt.Sprintf("/api/recon/vulnerability/list?targetid=%s", url.QueryEscape(targetID))
	resp, err := s.SAC.AuthGet(endpoint, token)
	if err != nil {
		return nil, err
	}
	dataField, ok := resp["data"]
	if !ok {
		return nil, fmt.Errorf("API response(ListVulnerabilities) missing 'data' field")
	}
	var vulns []png_hnd.Vulnerabilities
	if err := s.SAC.Decode(dataField, &vulns); err != nil {
		return nil, fmt.Errorf("failed to decode vulnerabilities: %w", err)
	}
    //fmt.Println(vulns)
	return vulns, nil
}