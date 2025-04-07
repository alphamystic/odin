package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Ratings struct {
  PluginName string
  Hash string // acts like the ID
  Rates float32
  Average int
  utils.TimeStamps
}

type Rater struct {
  UserID string
  PluginID string
  Comment string
  utils.TimeStamps
}
