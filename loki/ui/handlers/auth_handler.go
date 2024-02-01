package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
  //"github.com/alphamystic/odin/loki/lib/workers"
  dfn"github.com/alphamystic/odin/lib/definers"
  srvs"github.com/alphamystic/odin/loki/internal/services"
)


func (hnd *Handler) IsAuthenticated(req *http.Request) bool{
  session,_ := hnd.Store.Get(req,"session")
  _,ok := session.Values["token"].(string)
  if !ok {
    return false
  }
  return true
}

func (hnd *Handler) Register(res http.ResponseWriter, req *http.Request){
  if !Registration {
    hnd.TPl.ExecuteTemplate(res,"register.html","Sorry no user registration allowed for now. contact the admin. :)")
    return
  }
  if _,err := hnd.GetUDFromToken(req); err != nil {
    if errors.Is(err,dfn.UserNotLoggedIn){
      goto END
    }
    if errors.Is(err,dfn.NoCLaims){
      http.Redirect(res,req,"/",http.StatusSeeOther)
      return
    }
    goto END
  }
  END:
  if req.Method == "POST"{
    fmt.Println("Runing post register...")
    req.ParseForm()
    id := utils.GenerateUUID()
    name := req.FormValue("name")
    mail := req.FormValue("email")
    pass := req.FormValue("password")
    pass2 := req.FormValue("cpass")
    if !utils.CheckifStringIsEmpty(name){
      hnd.Tpl.ExecuteTemplate(res,"register.html","User name can not be empty")
      return
    }
    if !utils.CheckifStringIsEmpty(mail){
      hnd.Tpl.ExecuteTemplate(res,"register.html","Email can not be empty.")
      return
    }
    if !utils.CheckifStringIsEmpty(pass){
      hnd.Tpl.ExecuteTemplate(res,"register.html","Password can not be empty")
      return
    }
    if !utils.CheckifStringIsEmpty(pass2){
      hnd.Tpl.ExecuteTemplate(res,"register.html","Password confirmation can not be empty")
      return
    }
    if pass != pass2 {
      hnd.Tpl.ExecuteTemplate(res,"register.html","Passwords are varying.")
      return
    }
    pass,err := utils.HashPassword(pass)
    if err != nil{
      hnd.Tpl.ExecuteTemplate(res,"error.html",ErrorRes{
        ErrorCode: 500,
        ErrorMessage: "We are experiencing internal server issues. Please try again later",
        PrevUrl: "/register",
      })
    }
    user := dfn.User{
      UserID: id,
      OwnerID: id,
      UserName: name,
      Email: mail,
      Password: pass,
      Active: true,
      Anonymous:  false,//utils.Md5Hash( id + utils.RandNoLetter(5))
      Verify: true,
      Admin: true,
      CreatedAt: currentTime,
      UpdatedAt: currentTime,
    }
    userService := &srvs.UserDtataService{}
    err = userService.CreateUser(user)
    if err != nil{
      utils.Danger(fmt.Errorf("%q",err))
      tpl.ExecuteTemplate(res,"register.html","There was an internal error. Try again later.")
      return
    }
    http.Redirect(res,req,"/mkubwa",http.StatusSeeOther)
    return
  }
  if req.Method == "GET"{
    hnd.Tpl.ExecuteTemplate(res,"register.html",nil)
    return
  }
  http.Redirect(res,req,"/mkubwa",http.StatusSeeOther)
  return
}

//  @TODO Add a set expiry
func (hnd *Handler) Signin(res http.ResponseWriter, req *http.Request){
  if req.Method == "POST"{
    req.ParseForm()
    email :=  req.FormValue("mail")
    if !utils.IsValidEmail(email){
      hnd.Tpl.ExecuteTemplate(res,"login.html","Wrong Email provided.")
      return
    }
    pass := req.FormValue("password")
    if !utils.CheckifStringIsEmpty(pass){
      hnd.Tpl.ExecuteTemplate(res,"login.html","Password can not be empty.")
      return
    }
    user,err := services.Login(pass,email)
    if err != nil {
      utils.Logerror(err)
      if errors.Is(err,dfn.WrongPassword) {
        hnd.RL.LogRequestDetails(req, fmt.Sprintf("Wrong password attmept with email %s and password %s",email,pass))
        hnd.Tpl.ExecuteTemplate(res,"login.html","Wrong email or password provided.")
        return
      }
      hnd.Tpl.ExecuteTemplate(res,"error.html","We are experiencing internal server issues, please try again later. :)")
      return
    }
    ud := &UserData {
      UserID: ud.UserID,
      Admin: ud.Admin
    }
    token,err := hnd.GenerateJWT(ud)
    if err != nil {
      utils.Danger(err)
      errPage := ErrorPage {
        ErrorCode: 500,
        Message: "We are experiencing internal server issues, please try again later. :)",
        Back: "/mkubwa",
      }
      hnd.Tpl.ExecuteTemplate(res,"error.html",errPage)
      return
    }
    session,_ := hnd.Store.Get(req,"session")
    session.Values["token"] = token
    session.Save(req,res)
    //redirect to dashboard or get the dash data and execute dash
    http.Redirect(res,req,"/",http.StatusSeeOther)
    return
  }
  tpl.hnd.Tpl.ExecuteTemplate(res,"login.html",nil)
  return
}

func (hnd *Handler) Logout(res http.ResponseWriter, req *http.Request){
  session,_ := hnd.Store.Get(req,"session")
  delete(session.Values,"token")
  session.Save(req,res)
  hnd.Tpl.ExecuteTemplate(res,"login.html","Logged Out ADIOS!!!")
  return
}
