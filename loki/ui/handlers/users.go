package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
)

func (hnd *Handler) RegularUsers(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"users-regular.html",nil)
  return
}

func (hnd *Handler) Admins(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"users-admins.html",nil)
  return
}
