package services

import (
  dfn"github.com/alphamystic/odin/lib/definers"
)

type (
  UserData interface {
    CreateUser() error
    ViewUsers()
    ListUsers()
    DeleteUser()
    IsAdmin()
    UpdateUser() error
  }
  UserDtataService struct {}
)

func (usd *UserDtataService) CreateUser(user dfn.User) error {
  return nil
}



func ViewUser(userId string)(*User,error){
}
