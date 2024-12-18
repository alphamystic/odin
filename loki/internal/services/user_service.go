package services

import (
  dfn"github.com/alphamystic/odin/lib/definers"
)

type (
  UserData interface {
    CreateUser(dfn.User) error
    ViewUsers()
    ListUsers()
    DeleteUser()
    IsAdmin()
    UpdateUser() error
  }
  UserDataService struct {}
)

func (usd *UserDataService) CreateUser(user dfn.User) error {
  return nil
}

func NewUserService() *UserDataService{
  return &UserDataService{}
}


func ViewUser(userId string)(*dfn.User,error){
  return nil,nil
}
