package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) ListAssets(res http.ResponseWriter, req *http.Request) {
	ud, authenticated := hnd.AuthenticateUser(res, req)
	if !authenticated {
		return
	}

	token, _ := hnd.GetToken(req)
	active := req.URL.Query().Get("active") != "false"
	hardware := req.URL.Query().Get("hardware") == "true"

	endpoint := fmt.Sprintf("/api/assets/list?hardware=%v&active=%v", hardware, active)
	resp, err := hnd.SRVCS.AssetSrvs.SAC.AuthGet(endpoint, token)

	// Initialize empty slice
	var renderAssets []map[string]interface{}

	// Check if API call failed or data is missing
	if err != nil || resp["data"] == nil {
		utils.Warning(fmt.Sprintf("Assets empty or error: %v", err))
		// ADD DUMMY ASSET TO PREVENT CRASH AND SHOW UI
		renderAssets = append(renderAssets, map[string]interface{}{
			"AssetID":          "0000-DUMMY",
			"AssetName":        "No Assets Found",
			"Description":      "Connect an agent or create an asset to see data here.",
			"Active":           false,
			"Hardware":         false,
			"DescMap":          map[string]string{"Status": "Empty"},
			"RawDescriberJSON": "{}",
		})
	} else {
		dataField, ok := resp["data"].([]interface{})
		if ok {
			for _, item := range dataField {
				asset, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				parsedDesc := make(map[string]string)
				if rawDesc, ok := asset["describer"].(string); ok {
					trimmed := strings.TrimPrefix(rawDesc, "map[")
					trimmed = strings.TrimSuffix(trimmed, "]")
					pairs := strings.Split(trimmed, " ")
					for _, p := range pairs {
						kv := strings.Split(p, ":")
						if len(kv) == 2 {
							parsedDesc[kv[0]] = kv[1]
						}
					}
				}

				descJSON, _ := json.Marshal(parsedDesc)

				renderAssets = append(renderAssets, map[string]interface{}{
					"AssetID":          asset["asset_id"],
					"AssetName":        asset["asset_name"], // Fixed: removed stray 'sla'
					"Description":      asset["description"],
					"Active":           asset["active"],
					"Hardware":         asset["hardware"],
					"DescMap":          parsedDesc,
					"RawDescriberJSON": string(descJSON),
				})
			}
		}
	}

	// ALWAYS EXECUTE TEMPLATE
	tpl, err := hnd.Pages.GetATemplate("asset_list", "asset_list.tmpl")
	if err != nil {
		utils.Warning(fmt.Sprintf("Tpl Load Error: %v", err))
		return
	}

	if err := tpl.ExecuteTemplate(res, "asset_list", map[string]interface{}{
		"UserData": ud,
		"assets":   renderAssets,
	}); err != nil {
		utils.Warning(fmt.Sprintf("Template Execution Error: %v", err))
	}
}