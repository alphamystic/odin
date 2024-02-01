package services

type (
  Authorize interface {
    Login() (string,error)
    Logout(string) error
  }
  AuthorizeService struct{
    //authDomain domain.
  }
)

func NewAuthorizeService() *AuthorizeService{
  return &AuthorizeService{}
}


func (as *AuthorizeService) Login()(string,error){
  return "",nil
}

func (as *AuthorizeService) ChangePassword() error {
  return nil
}

func (as *AuthorizeService) Logout(token string) error  {
  return nil
}
