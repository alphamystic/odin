package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
)


func Blank(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"blank.html",nil)
  return
}

func Home(res http.ResponseWriter, req *http.Request){
  /*if !IsAuthenticated(req){
    http.Redirect(res,req,"/mkubwa",http.StatusFound)//302
    return
  }*/
  if req.Method == "GET"{
    hnd.Tpl.ExecuteTemplate(res,"home.html",nil)
    return
  }
  http.Redirect(res,req,"/logout",http.StatusFound)
  return
}
