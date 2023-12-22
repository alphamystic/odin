package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
  srvs"github.com/alphamystic/odin/loki/internal/services"
)

func (hnd *Handler) Listapikeys(res http.ResponseWriter, req *http.Request){
  if req.Method != "GET" {
    errPage := ErrorPage {
      ErrorCode: 400,
      Message: "Get a life dummy, anyway if you find a vulnerablity fix it or email us vulnerablity.odin@eloracle.africa",
      Back:"/logout",
    }
    hnd.Tpl.ExecuteTemplate(res,"error.html",errPage)
    return
  }
  ntfs,err := srvs.ListUserNotifications(ud.Uid)
  if err != nil {
    utils.Warning(fmt.Sprintf("%s",err))
  }
  hnd.Tpl.ExecuteTemplate(res,"list-apikey.html",LOKI{
    "notifications": ntfs,
  })
  return
}

func (hnd *Handler) Createapikeys(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    ntfs,err := srvs.ListUserNotifications(ud.Uid)
    if err != nil {
      utils.Warning(fmt.Sprintf("%s",err))
    }
    hnd.Tpl.ExecuteTemplate(res,"create-apikey.html",LOKI{
      "notifications": ntfs,
    })
  }
  ntfs,err := srvs.ListUserNotifications(ud.Uid)
  if err != nil {
    utils.Warning(fmt.Sprintf("%s",err))
  }
  hnd.Tpl.ExecuteTemplate(res,"create-apikey.html",LOKI{
    "notifications": ntfs,
  })
  return
}




func GenerateToken(rt *Runtime) (string,error){
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

func AuthMiddleware(next http.Handler) http.Handler{
  return http.HandlerFunc(func(res http.ResponseWriter,req *http.Request){
    //tokenString := req.Header.Get("Authorization")
    /*token,err := jwt.ParseWithClaims(tokenString, &Runtime{},func(token *jwt.Token) (interface{},error){
      return store,nil
    })*/
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
        Code: http.StatusUnauthorized,
        Data: "Login to interact with server.",
        Message: "You are not logged in. Please log in first",
        Back: "Login.",
        Direction: "/login",
      }
      tpl.ExecuteTemplate(res,"error.html",tmplError)
      return
    }
    next.ServeHTTP(res,req)
  })
}
