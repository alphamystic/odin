package handlers

import (
  //"fmt"
  dom"github.com/alphamystic/odin/lib/domain"
  srvs"github.com/alphamystic/odin/loki/internal/services"
)
type Services struct {
  NTFCNSvrs *srvs.NotificationService
  AuthSrvs *srvs.AuthorizeService
  UserSrvs *srvs.UserDataService
}

func InitializeServices(domain *dom.Domain) *Services {
  notificationService := srvs.CreateNotifyer(domain)
  auth_service := srvs.NewAuthorizeService(domain)
  user_service := srvs.NewUserService(domain)
  return &Services{
    NTFCNSvrs: notificationService,
    AuthSrvs: auth_service,
    UserSrvs: user_service,
  }
}
