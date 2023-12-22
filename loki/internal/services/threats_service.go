package services

/*
  * So a service basicaly calls the back end function interacting with the DB from the domain
*/
type (
  ThreatService interface {
    CreateThrreat()
    ViewThreat()
    ListThreat()
    MarkThreatAsActice()
  }
)
