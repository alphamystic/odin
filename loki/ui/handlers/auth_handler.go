package handlers

import(
  "os"
  "fmt"
  "io"
  "errors"
  "net/http"
  "path/filepath"
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
    hnd.ExecRegister(res,req, "Sorry no user registration allowed for now. contact the admin. :)")
    return
  }
  //END:
  if req.Method == "POST"{
    err := req.ParseMultipartForm(10 << 1) // 10MB limit
		if err != nil {
			utils.Warning(fmt.Sprintf("Error parsing login multipart form: %s", err))
			hnd.Internalserverror(res, req)
			return
		}
    req.ParseForm()
    id := utils.GenerateUUID()
    name := req.FormValue("name")
    mail := req.FormValue("email")
    pass := req.FormValue("password")
    pass2 := req.FormValue("cpass")
    if !utils.CheckifStringIsEmpty(name){
      hnd.ExecLogin(res, req,"User name can not be empty")
      return
    }
    if !utils.CheckifStringIsEmpty(mail){
      hnd.ExecRegister(res, req, "Email can not be empty.")
      return
    }
    if !utils.CheckifStringIsEmpty(pass){
      hnd.ExecRegister(res, req, "Password can not be empty")
      return
    }
    if !utils.CheckifStringIsEmpty(pass2){
      hnd.ExecRegister(res, req, "Password confirmation can not be empty")
      return
    }
    if pass != pass2 {
      hnd.ExecRegister(res, req,"Passwords are varying.")
      return
    }
    pass,err = utils.HashPassword(pass)
    if err != nil{
      hnd.ExecRegister(res, req, "We are experiencing internal server issues. Please try again later")
    }
    ctx := req.Context()
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
    err = hnd.SRVCS.UserSrvs.CreateUser(ctx,user)
    if err != nil{
      utils.Danger(fmt.Errorf("%q",err))
      hnd.ExecRegister(res,req,"There was an internal error. Try again later.")
      return
    }
    // save the image profile
    dir := fmt.Sprintf("./loki/ui/static/uploads/profile/%s", user.UserID)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			utils.Warning(fmt.Sprintf("Error creating directory: %s", err))
      // changing from internal server error to  login but with update profile picture at the bottom
			hnd.ExecLogin(res, req,"Please update your profile picture on login in.")
			return
		}
    files := req.MultipartForm.File["images"]
		for _, fileHeader := range files {
			// Open the uploaded file
			file, err := fileHeader.Open()
			if err != nil {
				utils.Warning(fmt.Sprintf("Error opening file: %s", err))
				hnd.ExecLogin(res, req,"Please update your profile picture on login in.")
				return
			}
			defer file.Close()
			// Create a temporary file to save the uploaded file
			tempFile, err := os.CreateTemp(dir, fmt.Sprintf("profile-%s-*%s", utils.GenerateUUID(), filepath.Ext(fileHeader.Filename)))
			if err != nil {
				utils.Warning(fmt.Sprintf("Error creating temporary file: %s", err))
				hnd.ExecLogin(res, req,"Please update your profile picture on login in.")
				return
			}
			defer tempFile.Close()
			// Copy the uploaded file to the temporary file
			_, err = io.Copy(tempFile, file)
			if err != nil {
				utils.Warning(fmt.Sprintf("Error saving file: %s", err))
				hnd.ExecLogin(res, req,"Please update your profile picture on login in.")
				return
			}
			break
		}
    http.Redirect(res,req,"/mkubwa",http.StatusSeeOther)
    return
  }
  if req.Method == "GET"{
    hnd.ExecRegister(res,req,"")
    return
  }
  http.Redirect(res,req,"/mkubwa",http.StatusSeeOther)
  return
}

/*
// You can try picking the user data and log it out for testing/logging purposes
if _,err := hnd.GetUDFromToken(req); err != nil {
  utils.Warning(fmt.Sprintf("%s", err))
  if errors.Is(err,dfn.UserNotLoggedIn){
    http.Redirect(res,req,"/register",http.StatusSeeOther)
    return
  }
  if errors.Is(err,dfn.NoClaimsError){
    http.Redirect(res,req,"/logout",http.StatusSeeOther)
    return
  }
  //we can log this error somewhere as a http request error
  hnd.ExecRegister(res, req,"")
  return
}*/
//  @TODO Add a set expiry
func (hnd *Handler) Signin(res http.ResponseWriter, req *http.Request){
  if req.Method == "POST"{
    ctx := req.Context()
    req.ParseForm()
    email :=  req.FormValue("mail")
    if !utils.IsValidEmail(email){
      hnd.ExecLogin(res,req,"Wrong Email provided.")
      return
    }
    pass := req.FormValue("password")
    if !utils.CheckifStringIsEmpty(pass){
      hnd.ExecLogin(res,req,"Password can not be empty.")
      return
    }
    user,err := hnd.SRVCS.AuthSrvs.Login(ctx,pass,email)
    if err != nil {
      utils.Logerror(err)
      if errors.Is(err,dfn.WrongPassword) {
        hnd.RL.LogRequestDetails(req, fmt.Sprintf("Wrong password attmept with email %s and password %s",email,pass))
        hnd.ExecLogin(res,req,"Wrong email or password provided.")
        return
      }
      hnd.ExecLogin(res,req,"We are experiencing internal server issues, please try again later. :)")
      return
    }
    ud := &UserData {
      UserId: user.UserID,
      UserName: user.UserName,
      Admin: user.Admin,
    }
    token,err := hnd.GenerateJWT(ud)
    if err != nil {
      utils.Danger(err)
      hnd.ExecLogin(res,req,"We are experiencing internal server issues, please try again later. :)")
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
  hnd.ExecLogin(res,req,"")
  return
}

// Find a way to store the cookie
func (hnd *Handler) Logout(res http.ResponseWriter, req *http.Request){
  _,err := req.Cookie("Authorization")
  if err == http.ErrNoCookie {
    hnd.ExecLogin(res,req,"Please login first.")
    return
  } else if err != nil {
      fmt.Println("[+]  Some internal error. \nERROR: ",err)
      hnd.ExecLogin(res,req,"Internal error, try again later.")
      return
  }
  //tokenString := cookie.Value
  req.Header.Del("Authorization")
  res.Header().Del("Authorization")
  //InvalidTokens = append(InvalidTokens,tokenString)
  hnd.ExecLogin(res,req,"Logged Out. ADIOS!!!")
  return
}

func isImageFile(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	}
	return false
}

func (hnd *Handler) ExecLogin(res http.ResponseWriter, req *http.Request,data string){
  tpl,err := hnd.Pages.GetAStaticTemplate("login","login.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"login",data)
  return
}

func (hnd *Handler) ExecRegister(res http.ResponseWriter, req *http.Request,data string){
  tpl,err := hnd.Pages.GetAStaticTemplate("register","register.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"register",data)
  return
}
