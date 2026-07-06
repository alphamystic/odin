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
  AssetSrvs *srvs.AssetService
  ReportSrvs *srvs.ReportService
  ApiKeySrvs *srvs.ApikeyService
  IssAppSrvs  *srvs.IssuesAppointmentService
  MinionSrvs *srvs.MinionService
  MothershipSrvs *srvs.MothershipService
  ReconSrvs *srvs.ReconService
  RMMSrvs    *srvs.RMMService
}

// change the services instead odf taking in a domain to take in
// a connection to the API.
func InitializeServices(baseURL, apiKey string) *Services {
  sac := srvs.NewServerAPIConnector(baseURL, apiKey)
  notificationService := srvs.CreateNotifyer(sac)
  auth_service := srvs.NewAuthorizeService(sac)
  user_service := srvs.NewUserService(sac)
  asset_service := srvs.NewAssetService(sac)
  report_service := srvs.NewReportService(sac)
  apikey_service := srvs.NewApikeyService(sac)
  minion_service := srvs.NewMinionService(sac)
  mothership_service := srvs.NewMothershipService(sac)
  recon_service := srvs.NewReconService(sac)
  rmm_service := srvs.NewRMMService(sac)
  issues_appointments_service := srvs.NewIAService(sac)
  return &Services{
    NTFCNSvrs: notificationService,
    AuthSrvs: auth_service,
    UserSrvs: user_service,
    AssetSrvs: asset_service,
    ReportSrvs: report_service,
    ApiKeySrvs: apikey_service,
    IssAppSrvs: issues_appointments_service,
    MinionSrvs: minion_service,
    MothershipSrvs: mothership_service,
    ReconSrvs:  recon_service,
    RMMSrvs:    rmm_service,
  }
}
