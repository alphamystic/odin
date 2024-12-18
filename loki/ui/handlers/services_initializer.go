package handlers

import (
  //"fmt"
  srvs"github.com/alphamystic/odin/loki/internal/services"
)
type Services struct {
  NTFCNSvrs *srvs.NotificationService
  AuthSrvs *srvs.AuthorizeService
  UserSrvs *srvs.UserDataService
}

func InitializeServices() *Services {
  notificationService := srvs.CreateNotifyer()
  auth_service := srvs.NewAuthorizeService()
  user_service := srvs.NewUserService()
  return &Services{
    NTFCNSvrs: notificationService,
    AuthSrvs: auth_service,
    UserSrvs: user_service,
  }
}
