package services

import(
  "context"
  dom"github.com/alphamystic/odin/lib/domain"
  dfn"github.com/alphamystic/odin/lib/definers"
)

type (
  Authorize interface {
    Login() (*dfn.User,error)
    Logout(string) error
  }
  AuthorizeService struct{
    //authDomain domain.
    Dom *dom.Domain
  }
)

func NewAuthorizeService(domain *dom.Domain) *AuthorizeService{
  return &AuthorizeService{Dom: domain}
}


func (as *AuthorizeService) Login(ctx context.Context,pass,email string)(*dfn.User,error){
  return as.Dom.Authenticate(ctx,pass,email)
}

func (as *AuthorizeService) ChangePassword() error {
  return nil
}

func (as *AuthorizeService) Logout(token string) error  {
  return nil
}
