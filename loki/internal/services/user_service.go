package services


import (
  "fmt"
  "context"
	//"encoding/json"
  dfn"github.com/alphamystic/odin/lib/definers"
)

type (
  UserData interface {
    CreateUser(ctx context.Context, u dfn.User) error
    ViewUsers(ctx context.Context)
    ListUsers(ctx context.Context)
    DeleteUser(ctx context.Context)
    IsAdmin(ctx context.Context)
    UpdateUser(ctx context.Context) error
  }
  UserDataService struct {SAC *ServerAPIConnector}
)

func NewUserService(sac *ServerAPIConnector) *UserDataService{
  return &UserDataService{SAC: sac}
}

func (usd *UserDataService) CreateUser(ctx context.Context,user dfn.User) (string,error) {
  data,err := usd.SAC.Post("/api/users/createuser",user)
  if err != nil{
    return  "",err
  }
  if status, ok := data["Status"].(string); ok {
		return "",fmt.Errorf("Status: %q", status)
	}
  redirecturl, ok := data["RedirectUrl"].(string)
  if !ok {
		return "",fmt.Errorf("Redirect URL: %q", redirecturl)
	}
  return redirecturl,nil
}


func (usd *UserDataService) ViewUser(userId string)(*dfn.User,error){
  return nil,nil
}
