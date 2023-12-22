package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Activity struct{
  //ThreatId string `json:"threatid"`
  Name string `json:"activity_name"`
  AId string  `json:"activity_id"`
  CreatorId string `json:"creator_id"`
  Description string  `json:"activity_description"`
  Validated bool `json:"validated"`
  utils.TimeStamps
}
