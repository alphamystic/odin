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

func (hnd *Handler) Register(res http.ResponseWriter, req *http.Request) {
  if !Registration {
    hnd.ExecRegister(res, req, "Sorry, user registration is disabled. Contact admin.")
    return
  }

  if req.Method == "POST" {
    err := req.ParseMultipartForm(10 << 20) // 10MB limit
    if err != nil {
      utils.Warning(fmt.Sprintf("Error parsing multipart form: %s", err))
      hnd.ExecRegister(res, req, "Invalid form submission. Ensure file uploads are enabled.")
      return
    }
    req.ParseForm()
    id := utils.GenerateUUID()
    name := req.FormValue("name")
    mail := req.FormValue("email")
    pass := req.FormValue("password")
    pass2 := req.FormValue("cpass")
    if !utils.CheckifStringIsEmpty(name) {
      hnd.ExecRegister(res, req, "Username cannot be empty")
      return
    }
    if !utils.CheckifStringIsEmpty(mail) {
      hnd.ExecRegister(res, req, "Email cannot be empty.")
      return
    }
    if !utils.CheckifStringIsEmpty(pass) {
      hnd.ExecRegister(res, req, "Password cannot be empty")
      return
    }
    if pass != pass2 {
      hnd.ExecRegister(res, req, "Passwords do not match.")
      return
    }
    ctx := req.Context()
    user := dfn.User{
      UserID:   id,
      OwnerID:  id,
      UserName: name,
      Email:    mail,
      Password: pass,
      Active:   true,
      Anonymous: false,
      Verify:   true,
      Admin:    true,
    }
    userID, err := hnd.SRVCS.UserSrvs.CreateUser(ctx, user)
    if err != nil {
      utils.Danger(fmt.Errorf("%q", err))
      hnd.ExecRegister(res, req, "Internal error. Try again later.")
      return
    }
    // Handle file upload
    dir := fmt.Sprintf("./loki/ui/static/uploads/profile/%s", userID)
    if err = os.MkdirAll(dir, os.ModePerm); err != nil {
      utils.Warning(fmt.Sprintf("Error creating directory: %s", err))
      hnd.ExecLogin(res, req, "Profile created. Please update your profile picture after logging in.")
      return
    }
    files := req.MultipartForm.File["images"]
    if len(files) > 0 {
      fileHeader := files[0] // Only save the first image
      file, err := fileHeader.Open()
      if err != nil {
        utils.Warning(fmt.Sprintf("Error opening file: %s", err))
        hnd.ExecLogin(res, req, "Profile created. Please update your profile picture after logging in.")
        return
      }
      defer file.Close()
      filePath := fmt.Sprintf("%s/profile-%s%s", dir, userID, filepath.Ext(fileHeader.Filename))
      outFile, err := os.Create(filePath)
      if err != nil {
        utils.Warning(fmt.Sprintf("Error creating file: %s", err))
        hnd.ExecLogin(res, req, "Profile created. Please update your profile picture after logging in.")
        return
      }
      defer outFile.Close()
      if _, err = io.Copy(outFile, file); err != nil {
        utils.Warning(fmt.Sprintf("Error saving file: %s", err))
        hnd.ExecLogin(res, req, "Profile created. Please update your profile picture after logging in.")
        return
      }
    }
    http.Redirect(res, req, "/mkubwa", http.StatusSeeOther)
    return
  }
  if req.Method == "GET" {
    hnd.ExecRegister(res, req, "")
    return
  }
  http.Redirect(res, req, "/mkubwa", http.StatusSeeOther)
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

    token,err := hnd.SRVCS.AuthSrvs.Login(ctx,pass,email)
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
    cookie := http.Cookie{
        Name:     "Authorization",
        Value:    token,
        Path:     "/",
        MaxAge:   72000,
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
