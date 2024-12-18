package handlers

import(
  "fmt"
  "errors"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
  //"github.com/alphamystic/odin/loki/lib/workers"
  dfn"github.com/alphamystic/odin/lib/definers"
)


func (hnd *Handler) IsAuthenticated(req *http.Request) bool{
  _,err := hnd.GetUDFromToken(req)
  if err != nil{
    return false
  }
  return true
}

func (hnd *Handler) Register(res http.ResponseWriter, req *http.Request){
  if !Registration {
    hnd.Tpl.ExecuteTemplate(res,"register.html","Sorry no user registration allowed for now. contact the admin. :)")
    return
  }
  if _,err := hnd.GetUDFromToken(req); err != nil {
    if errors.Is(err,dfn.UserNotLoggedIn){
      goto END
    }
    if errors.Is(err,dfn.NoClaimsError){
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
        "ErrorCode": 500,
        "Data": "We are experiencing internal server issues. Please try again later",
        "Direction": "/register",
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
    }
    user.Touch()
    err = hnd.SRVCS.UserSrvs.CreateUser(user)
    if err != nil{
      utils.Danger(fmt.Errorf("%q",err))
      hnd.Tpl.ExecuteTemplate(res,"register.html","There was an internal error. Try again later.")
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
    user,err := hnd.SRVCS.AuthSrvs.Login(pass,email)
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
      UserId: user.UserID,
      Admin: user.Admin,
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
    cookie := http.Cookie{
        Name:     "Authorization",
        Value:    token,
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(res,&cookie)
    //redirect to dashboard or get the dash data and execute dash
    http.Redirect(res,req,"/",http.StatusSeeOther)
    return
  }
  hnd.Tpl.ExecuteTemplate(res,"login.html",nil)
  return
}

// Find a way to store the cookie
func (hnd *Handler) Logout(res http.ResponseWriter, req *http.Request){
  _,err := req.Cookie("Authorization")
  if err == http.ErrNoCookie {
    hnd.Tpl.ExecuteTemplate(res,"login.html","Please login first.")
    return
  } else if err != nil {
      fmt.Println("[+]  Some internal error. \nERROR: ",err)
      hnd.Tpl.ExecuteTemplate(res,"login.html","Internal error, try again later.")
      return
  }
  //tokenString := cookie.Value
  req.Header.Del("Authorization")
  res.Header().Del("Authorization")
  //InvalidTokens = append(InvalidTokens,tokenString)
  hnd.Tpl.ExecuteTemplate(res,"signin.html","Logged Out. ADIOS!!!")
  return
}
