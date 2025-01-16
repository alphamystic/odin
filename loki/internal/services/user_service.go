package services


import (
  "context"
  dom"github.com/alphamystic/odin/lib/domain"
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
  UserDataService struct {Dom *dom.Domain}
)

func NewUserService(domain *dom.Domain) *UserDataService{
  return &UserDataService{Dom: domain}
}

func (usd *UserDataService) CreateUser(ctx context.Context,user dfn.User) error {
  return usd.Dom.CreateUser(ctx, user)
}


func (usd *UserDataService) ViewUser(userId string)(*dfn.User,error){
  return nil,nil
}
