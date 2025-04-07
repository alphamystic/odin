package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Asset struct{
  Name string `json:"asset_name"`
  AssetID string  `json:"asset_id"`
  Description string  `json:"description"`
  Dscbr interface{} `json:"describer"`
  Active bool `json:"active"`
  Hard bool `json:"hardware"`
  utils.TimeStamps
}

// this is marshalled up from a jwt token
type Describer struct {
  AgentID string `json:"agentid"`
  UserId string  `json:"uid"`
  GroupId string  `json:"groupid"`
  HomeDir string `json:"homedir"`
  Ops  string  `json:"ostype"`
  Description string `json:"description"`
  Installed bool `json:"installed"`
  AType string `json:"asset_type"` //.dll .exe .so
}
