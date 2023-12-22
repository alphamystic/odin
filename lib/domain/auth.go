package domain

/*
  This package defines authorizing functions
*/
import (
  "context"
  "database/sql"
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

const (
  authStmt = "SELECT userid,username,email,password,active FROM `odin`.`user` WHERE email = ?;"
)

func (d *Domain) Authenticate(password,email string,ctx context.Context)  (*dfn.User,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var userEmail,userName,hash,userId string
  var active,admin bool
  row :=  conn.QueryRowContext(ctx,authStmt,email)
  err := row.Scan(&userId,&userName,&userEmail,&hash,&active,&admin)
  if err != nil{
    if err == sql.ErrNoRows {
      utils.Warning(fmt.Sprintf("A none user with email %s tried accessing server",email))
      return &dfn.User{},err
    }
    _ = utils.LogToFile(utils.Logger{Name: "auth_sql", Text: fmt.Sprintf("A non-user with email %s tried accessing the server", email)})
    return &dfn.User{},err
  }
  if !active{
    _ = utils.LogToFile(utils.Logger{Name: "auth_danger_sql", Text: fmt.Sprintf("A non-active with email %s tried accessing the server", email),})
    return &dfn.User{},dfn.NonActiveUser
  }
  err = utils.CheckPasswordHash(password, hash)
  if err != nil{
    _ = utils.LogToFile(utils.Logger{Name: "auth_danger", Text: fmt.Sprintf("Wrong login attempt for email %s with password %s  %s", email, password, err),})
    utils.Warning(fmt.Sprintf("Wrong login attempt for email %s with password %s. ERROR: %s\n",email,password,err))
    return &dfn.User{},dfn.WrongPassword
  }
  return &dfn.User {
    UserID: userId,
    UserName: userName,
    Admin: admin,
    },nil
}
