package domain


import (
  "fmt"
  "errors"
  "context"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

// https://www.youtube.com/watch?v=iJTtSd_wJtQ
const (
  createHashStmt = `INSERT INTO odin.hashes (userid,hash,created_at,updated_at) VALUES(?,?,?,?);`
  getHash = "SELECT * FROM `odin`.`hashes` WHERE userid = ?;"
  updateHash = `UPDATE odin.hashes SET (hash = ? AND updated_at = ?) WHERE (hash = ? AND userid = ?);`
  createUser = `INSERT INTO odin.user (userid,ownerid,username,email,password,active,anonymous,verified,admin,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)`
  viewUser = `SELECT * FROM odin.user WHERE userid = ?;`
  listMyUsers = `SELECT * FROM odin.user WHERE (active = ? AND verified = ? AND ownerid = ?) ORDER BY updated_at ASC;`
  adminListUsers = `SELECT * FROM odin.user WHERE (active = ?) ORDER BY updated_at ASC;`
  listAllUsers =`SELECT * FROM odin.user WHERE (active = ?) ORDER BY updated_at ASC;`
  checkIfOwner = `SELECT userid,ownerid FROM odin.user WHERE (userid = ? AND ownerid = ?);`
  checkIfVerified = `SELECT userid,email,verified FROM odin.users WHERE (userid = ? AND email = ?);`
  checkIfAdmin = `SELECT active,admin FROM odin.user WHERE (userid= ?);`
)

type UserDom struct {
  Dbs *sql.DB
}
func (d *Domain) CreateHash(h dfn.UserHash,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %w",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err = conn.PrepareContext(ctx,updateHash)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error preparing to create user hash: %s"),})
    return errors.New("Server encountered an error while preparing to create hash. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&h.UserID,&h.Hash,&h.CreatedAt,&h.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error executing create user hash: %s",err),})
    return errors.New("Server encountered an error while creating hash.")
  }
  return nil
}

func (d *Domain) UpdateHash(prevHash,newHash,userId string,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  stmt,err := conn.PrepareContext(ctx,updateHash)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error preparing context to update hash: %s",err),})
    return errors.New("Server encountered an error while preparing to update hash")
  }
  defer stmt.Close()
  var res sql.Result
  res,err = stmt.ExecContext(ctx,newHash,utils.GetCurrentTime(),prevHash,userId)
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error executing update hash %s. ERROR: %s",prevHash,err),})
    return errors.New("Server encountered an error while executing update hash.")
  }
  return nil
}

func (d *Domain) GetHash(userId string,ctx context.Context) (*dfn.UserHash,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var h dfn.UserHash
  row := conn.QueryRowContext(ctx,getHash,userId)
  err = row.Scan(&h.UserID,&h.Hash,&h.CreatedAt,&h.UpdatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Hash for userid %s does not exist: ERROR: %",userId,err),})
      return nil,errors.New(fmt.Sprintf("Userid of %s is non existance.",userId))
    }
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error viewing user hash %s. ERROR: %s",userId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing user with id of %s",userId))
  }
  return &h,nil
}
// 	userid 	ownerid 	username 	email 	password 	active 	anonymous 	verified 	admin 	created_at 	updated_at
func (d *Domain) CreateUser(ctx context.Context,u dfn.User) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err = conn.PrepareContext(ctx,createUser)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error preparing to create user: %s",err),})
    return errors.New("Server encountered an error while preparing to create user. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&u.UserID,&u.OwnerID,&u.UserName,&u.Email,&u.Password,&u.Active,&u.Anonymous,&u.Verify,&u.Admin,&u.CreatedAt,&u.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error executing create user: %s",err),})
    return errors.New("Server encountered an error while creating user.")
  }
  return nil
}


func (d *Domain) ViewUser(userId string,ctx context.Context) (*dfn.User,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var u dfn.User
  row := conn.QueryRowContext(ctx,viewUser,userId)
  err = row.Scan(&u.UserID,&u.OwnerID,&u.UserName,&u.Email,&u.Password,&u.Active,&u.Anonymous,&u.Verify,&u.Admin,&u.CreatedAt,&u.UpdatedAt)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error viewing user %s. ERROR: %s",userId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing user with id of %s",userId))
  }
  return &u,nil
}

func (d *Domain) ListMyUsers(ctx context.Context,ownerId string,active,verified bool) ([]dfn.User,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
 rows,err := conn.QueryContext(ctx,listMyUsers,active,verified,ownerId)
 if err != nil{
   d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("ELAU: %s",err),})
   return nil,errors.New("Server encountered an error while listing all specified users.")
 }
 defer rows.Close()
 var users []dfn.User
 for rows.Next(){
   var u dfn.User
   err = rows.Scan(&u.UserID,&u.OwnerID,&u.UserName,&u.Email,&u.Password,&u.Active,&u.Anonymous,&u.Verify,&u.Admin,&u.CreatedAt,&u.UpdatedAt)
   if err != nil{
     d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error scanning list users: %s",err),})
     continue
   }
   users = append(users,u)
 }
 return users,nil
}


func (d *Domain) AdminListUsers(active bool,ctx context.Context) ([]dfn.User,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
 rows,err := conn.QueryContext(ctx,adminListUsers,active)
 if err != nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error listing admin users: %s",err),})
   //LogChecker()
   return nil,errors.New("Server encountered an error while listing all specified users.")
 }
 defer rows.Close()
 var users []dfn.User
 for rows.Next(){
   var u dfn.User
   err = rows.Scan(&u.UserID,&u.OwnerID,&u.UserName,&u.Email,&u.Password,&u.Active,&u.Anonymous,&u.Verify,&u.Admin,&u.CreatedAt,&u.UpdatedAt)
   if err != nil{
     d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error scanning admin list users: %s",err),})
     continue
   }
   users = append(users,u)
 }
 /* just return an empty list (even if an error occured it still will be empty)
 if len(users) < 0{
   return nil,errors.New(fmt.Sprintf("You probably have no users or an internal db issue, check logs for %s",utils.GetCurrentTime()))
 }*/
 return users,nil
}

func (d *Domain) ListAllUsers(active bool,ctx context.Context) ([]dfn.User,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
 rows,err := conn.QueryContext(ctx,listAllUsers,active)
 if err != nil{
   d.LogToFile(utils.Logger{Name:"users_sql",Text: fmt.Sprintf("Error listing all users: %w",err),})
   return nil,errors.New("Server encountered an error while listing all users.")
 }
 defer rows.Close()
 var users []dfn.User
 for rows.Next(){
   var u dfn.User
   err = rows.Scan(&u.UserID,&u.OwnerID,&u.UserName,&u.Email,&u.Password,&u.Active,&u.Anonymous,&u.Verify,&u.Admin,&u.CreatedAt,&u.UpdatedAt)
   if err != nil{
     d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error scaning for listed users: %w",err)})
     return nil,errors.New("Server encountered an error while listing allusers.")
   }
   users = append(users,u)
 }
 return users,nil
}

// only user and admin can deactivates user
func (d *Domain) MarkUserAsInActive(ownerId,userId string,active bool,ctx context.Context) error {
  // I complicated this for no good reason other than distracted thoughts
  var upStmt string
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  if utils.CheckifStringIsEmpty(ownerId) || ownerId == userId {
    upStmt = "UPDATE `odin`.`user` SET (`active` = ? AND `updated_at` = ?) WHERE (`userid` = ?);";
    goto PROCEED
  }
  if d.CheckIfAdmin(ownerId,ctx){
    // statement for admin
    upStmt = "UPDATE `odin`.`user` SET (`active` = ? AND `updated_at` = ?) WHERE (`userid` = ?);";
    goto PROCEED
  } else{
    return errors.New("A user can only be marked active or inactive by admin or the owner.")
  }
  defer conn.Close()
  //upStmt := "UPDATE `odin`.`users` SET (`active` = ? AND `updated_at` = ?) WHERE (`virusid` = ?);";
  PROCEED:
  stmt,err := conn.PrepareContext(ctx,upStmt)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error preparing to mark user %s as inactive: %s",err),})
    return errors.New("Server encountered an error while preparing to mark threat as active inactive.")
  }
  defer stmt.Close()
  var res sql.Result
  res,err = stmt.ExecContext(ctx,active,utils.GetCurrentTime(),userId)
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error marking user as active inactive: %s",err),})
    return errors.New("Server encountered an error while executing update apikey.")
  }
  return nil
}

// only admin or owner can verify
func (d *Domain) VerifyUser(ownerId,userId string,verified bool,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var upStmt string
  var admin bool
  if d.CheckIfAdmin(ownerId,ctx){
    // statement for admin
    admin = true
    upStmt = "UPDATE `odin`.`user` SET (`verified` = ? AND `updated_at` = ?) WHERE (`userid` = ?);";
  } else {
    admin = false
    upStmt = "UPDATE `odin`.`user` SET (`verified` = ? AND `updated_at` = ?) WHERE (`userid` = ? AND `ownerid` = ?);";
  }
  stmt,err := conn.PrepareContext(ctx,upStmt)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"sql",Text:fmt.Sprintf("EPTVU: %s",err)})
    utils.Logerror(err)
    return errors.New("Server encountered an error while preparing to verify user.")
  }
  defer stmt.Close()
  var res sql.Result
  if admin{
    res,err = stmt.Exec(verified,utils.GetCurrentTime(),userId)
  } else {
    res,err = stmt.Exec(verified,utils.GetCurrentTime(),userId,ownerId)
  }
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    d.LogToFile(utils.Logger{Name:"sql",Text:fmt.Sprintf("EEVU id: %sERROR: %s",userId,err),})
    utils.Logerror(err)
    return errors.New("Server encountered an error while executing update apikey.")
  }
  return nil
}

func (d *Domain) CheckIfOwner(userId,ownerId string,ctx context.Context) bool {
  var owid,uid string
  conn,err := d.GetConnection(ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error getting db connection: %q",err))
    return false
  }
  defer conn.Close()
  row := conn.QueryRowContext(ctx,checkIfOwner,userId,ownerId)
  err = row.Scan(&uid,&owid)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error scanning users rows %s",err)})
    return false
  }
  return true
}

//@ TODO: Check this put with a clear mind
func (d *Domain) CheckIfVerified(userId,email string,ctx context.Context) bool {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error getting db connection: %q",err))
    return false
  }
  var user string
  var verify,verified bool
  row := conn.QueryRowContext(ctx,checkIfVerified,userId,email,verify)
  err = row.Scan(&user,&verified)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"users_sql",Text:fmt.Sprintf("Error checking if verified %s",err),})
    return false
  }
  return true
}

func (d *Domain) CheckIfAdmin(adminId string,ctx context.Context) bool {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error getting db connection: %q",err))
    return false
  }
  var active,admin bool
  //"SELECT active,admin FROM `odin`.`users` WHERE (`userid`= ? AND `active` = ? AND `verified` = ?);"
  row := conn.QueryRowContext(ctx,checkIfAdmin,adminId)
  err = row.Scan(&active,&admin)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"danger_sql",Text:fmt.Sprintf("Error scanning for admin with id %s rows %s",adminId,err),})
    return false
  }
  if !active && admin{
    d.LogToFile(utils.Logger{Name:"danger_sql",Text:fmt.Sprintf("Old admin with id %s tried accessing admin stuff",adminId),})
    return false
  }
  if !active && !admin{
    d.LogToFile(utils.Logger{Name:"danger_sql",Text:fmt.Sprintf("Non admin with id %s tried accessing admin stuff",adminId)},)
    return false
  }
  return true
}
