package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) RegularUsers(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("users-regular","users-regular.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"users-regular",nil)
  return
}

func (hnd *Handler) Admins(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("users-admins","users-admins.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"users-admins",nil)
  return
}
