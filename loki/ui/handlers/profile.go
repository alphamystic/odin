package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) Profile(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("profile","profile.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
    return
  }
  tpl.ExecuteTemplate(res,"profile",nil)
  return
}

func (hnd *Handler) Updateprofile(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("profile-update","profile-update.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
    return
  }
  tpl.ExecuteTemplate(res,"profile-update",nil)
  return
}

func (hnd *Handler) Securityprofile(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("profile-security","profile-security.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
    return
  }
  tpl.ExecuteTemplate(res,"profile-security",nil)
  return
}

func (hnd *Handler) Notificationsprofile(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("profile-notifications","profile-notifications.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
    return
  }
  tpl.ExecuteTemplate(res,"profile-notifications",nil)
  return
}
