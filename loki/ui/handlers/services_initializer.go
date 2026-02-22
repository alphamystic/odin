package handlers

import (
  //"fmt"
  //dom"github.com/alphamystic/odin/lib/domain"
  srvs"github.com/alphamystic/odin/loki/internal/services"
)
type Services struct {
  NTFCNSvrs *srvs.NotificationService
  AuthSrvs *srvs.AuthorizeService
  UserSrvs *srvs.UserDataService
}

// change the services instead odf taking in a domain to take in
// a connection to the API.
func InitializeServices(baseURL, apiKey string) *Services {
  sac := srvs.NewServerAPIConnector(baseURL, apiKey)
  notificationService := srvs.CreateNotifyer(sac)
  auth_service := srvs.NewAuthorizeService(sac)
  user_service := srvs.NewUserService(sac)
  return &Services{
    NTFCNSvrs: notificationService,
    AuthSrvs: auth_service,
    UserSrvs: user_service,
  }
}
