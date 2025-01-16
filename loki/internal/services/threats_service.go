package services

/*
  * So a service basicaly calls the back end function interacting with the DB from the domain
*/

import (
  dom"github.com/alphamystic/odin/lib/domain"
)
type (
  ThreatService interface {
    CreateThrreat()
    ViewThreat()
    ListThreat()
    MarkThreatAsActice()
  }
  ThreatsService struct { Dom *dom.Domain }
)
