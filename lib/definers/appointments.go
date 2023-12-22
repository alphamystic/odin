package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Appointment struct {
  UserID string `json:"userid"`
  AppID string `json:"appointmentid"`
  Title string `json:"title"`
  Description string `json:"description"`
  Done bool `json:"done"`
  utils.TimeStamps
}
