package services

import(
  dfn"github.com/alphamystic/odin/lib/definers"
)

type (
  Authorize interface {
    Login() (*dfn.User,error)
    Logout(string) error
  }
  AuthorizeService struct{
    //authDomain domain.
  }
)

func NewAuthorizeService() *AuthorizeService{
  return &AuthorizeService{}
}


func (as *AuthorizeService) Login(pass,email string)(*dfn.User,error){
  return nil,nil
}

func (as *AuthorizeService) ChangePassword() error {
  return nil
}

func (as *AuthorizeService) Logout(token string) error  {
  return nil
}
