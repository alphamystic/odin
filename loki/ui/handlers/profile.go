package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
)

func (hnd *Handler) Profile(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"profile.html",nil)
  return
}

func (hnd *Handler) Updateprofile(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"profile-update.html",nil)
  return
}

func (hnd *Handler) Securityprofile(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"profile-security.html",nil)
  return
}

func (hnd *Handler) Notificationsprofile(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"profile-notifications.html",nil)
  return
}
