package services

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

}



func ViewUser(userId string)(*User,error){
}
