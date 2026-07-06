package services


import (
  "fmt"
  "context"
	//"encoding/json"
  dfn"github.com/alphamystic/odin/lib/definers"
)

// Reports are blogs which can be pentesting reports, Vulnerability/Pentesting and Audit Reports etc
type (
  Asset interface {
  	ListAssets(ctx context.Context, token string, active, hardware bool) ([]dfn.Asset, error)
  	//ListActiveAssets(ctx context.Context, token string, active bool, active bool) ([]dfn.Asset, error)
    ExportAsset(ctx context.Context) //Exporting a report
    ImportAsset(ctx context.Context)
  }
  AssetService struct {SAC *ServerAPIConnector}
)

func NewAssetService(sac *ServerAPIConnector) *AssetService{
  return &AssetService{SAC: sac}
}



// func (s *AssetService) ListActiveAssets(ctx context.Context, token string, active bool) ([]dfn.Asset, error) {
// 	endpoint := fmt.Sprintf("/api/assets/list?active=%s", active)
//     // 2. Make the authorized request
//     resp, err := s.SAC.AuthGet(endpoint, token)
//     if err != nil {
//         return nil, err
//     }
//     utils.Warning(fmt.Sprintf("%v",resp))
//     // 3. Extract the "data" field from the response map
//     // Your API returns: {"status": "success", "data": [...], "message": "..."}
//     dataField, ok := resp["data"]
//     if !ok {
//         return nil, fmt.Errorf("API response missing 'data' field")
//     }
//     utils.Warning(fmt.Sprintf("%v",resp))
//     var assets []dfn.Asset
//     // 4. Decode ONLY the data slice into your struct
//     if err := s.SAC.Decode(dataField, &parent_categories); err != nil {
//         return nil, fmt.Errorf("failed to decode Assets: %w", err)
//     }
//     utils.Warning(fmt.Sprintf("%v",assets))
//     return assets, nil
// }

func (s *AssetService) ListAssets(ctx context.Context, token string, active, hardware bool) ([]dfn.Asset, error) {
    endpoint := fmt.Sprintf("/api/assets/list?hardware=%v&active=%v", hardware, active)
    resp, err := s.SAC.AuthGet(endpoint, token)
    if err != nil {
        return nil, err
    }

    dataField, ok := resp["data"]
    if !ok {
        return nil, fmt.Errorf("API response missing 'data' field")
    }

    var assets []dfn.Asset
    // FIX: Decode into &assets, NOT &parent_categories
    if err := s.SAC.Decode(dataField, &assets); err != nil {
        return nil, fmt.Errorf("failed to decode Assets: %w", err)
    }
    return assets, nil
}
