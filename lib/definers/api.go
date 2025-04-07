package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Api struct {
  ApiKey string `json:"apikey"`
  OwnerID string `json:"ownerid"` // refers to apikey owner
  Active bool `json:"active"`
  utils.TimeStamps
}
