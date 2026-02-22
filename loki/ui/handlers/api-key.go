package handlers

import(
  "fmt"
  "time"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
//  "loki/lib/workers"
//  dfn"github.com/alphamystic/odin/lib/definers"

	"github.com/dgrijalva/jwt-go"
)

func (hnd *Handler) Listapikeys(res http.ResponseWriter, req *http.Request){
  if req.Method != "GET" {
    errPage := ErrorPage {
      ErrorCode: 400,
      Message: "Get a life dummy, anyway if you find a vulnerablity fix it or email us vulnerablity.odin@eloracle.africa",
      Back:"/logout",
    }
    tpl,err := hnd.Pages.GetATemplate("tmpl_error","templated_error.tmpl")
    if err != nil {
      utils.Warning(fmt.Sprintf("%s", err))
      hnd.Internalserverror(res, req)
  		return
    }
    tpl.ExecuteTemplate(res,"tmpl_error",errPage)
    return
  }
  ud, authenticated := hnd.AuthenticateUser(res, req)
  if !authenticated {
    return // User is redirected in the helper
  }
  ntfs,err := hnd.SRVCS.NTFCNSvrs.ListUserNotifications(ud.UserId)
  if err != nil {
    utils.Warning(fmt.Sprintf("%s",err))
  }

  tpl,err := hnd.Pages.GetATemplate("listapikey","listapikey.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"listapikey",LOKI{
    "notifications":ntfs,
    "userdata": ud,
  })
  return
}

func (hnd *Handler) Createapikeys(res http.ResponseWriter, req *http.Request){
  ud, authenticated := hnd.AuthenticateUser(res, req)
  if !authenticated {
    return // User is redirected in the helper
  }
  if req.Method != "POST"{
    ntfs,err := hnd.SRVCS.NTFCNSvrs.ListUserNotifications(ud.UserId)
    if err != nil {
      utils.Warning(fmt.Sprintf("%s",err))
    }
    tpl,err := hnd.Pages.GetATemplate("createapikey","createapikey.tmpl")
    if err != nil {
      utils.Warning(fmt.Sprintf("%s", err))
      hnd.Internalserverror(res, req)
  		return
    }
    tpl.ExecuteTemplate(res,"createapikey",LOKI{
      "Notifications": ntfs,
      "userdata": ud,
    })
  }
  tpl,err := hnd.Pages.GetATemplate("createapikey","createapikey.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  ntfs,err := hnd.SRVCS.NTFCNSvrs.ListUserNotifications(ud.UserId)
  if err != nil {
    utils.Warning(fmt.Sprintf("%s",err))
  }
  tpl.ExecuteTemplate(res,"create-apikey.html",LOKI{
    "Notifications": ntfs,
    "userdata": ud,
  })
  return
}




func GenerateToken(rt *UserData) (string,error){
  expTime := time.Now().Add(time.Hour * 48)
  /*rt.StandardClaims = jwt.StandardClaims{
    ExpiresAt: expTime.Unix(),
  }*/
  token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
    "runtime": rt,
    "exp": expTime.Unix(),
  })
  sighnedToken,err := token.SignedString(store)
  if err != nil{
    return "",err
  }
  return sighnedToken,nil
}



/*
func (hnd *Handler) AuthMiddleware(next http.Handler) http.Handler{
  return http.HandlerFunc(func(res http.ResponseWriter,req *http.Request){
    //tokenString := req.Header.Get("Authorization")
    /*token,err := jwt.ParseWithClaims(tokenString, &Runtime{},func(token *jwt.Token) (interface{},error){
      return store,nil
    })
    cookie,_ := req.Cookie("Authorization")
    tokenString := cookie.Value
    token,err := jwt.Parse(tokenString,func(tkn *jwt.Token)(interface{},error){
      if tkn.Method != jwt.SigningMethodHS256 {
        return nil,fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"])
      }
      return store,nil
    })
    if err != nil || !token.Valid {
      fmt.Println("[-]ERROR: during  authneticating token. \n  ",err)
      res.WriteHeader(http.StatusUnauthorized)
      tmplError := ErrorPage {
        ErrorCode: http.StatusUnauthorized,
        Data: "Login to interact with server.",
        Message: "You are not logged in. Please log in first",
        Back: "Login.",
        Direction: "/mkubwa",
      }
      tpl.ExecuteTemplate(res,"error.html",tmplError)
      return
    }
    next.ServeHTTP(res,req)
  })
}
*/
