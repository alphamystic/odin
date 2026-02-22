package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) Bugbounty(res http.ResponseWriter, req *http.Request){
  _, authenticated := hnd.AuthenticateUser(res, req)
  if !authenticated {
    return // User is redirected in the helper
  }
  tpl,err := hnd.Pages.GetATemplate("blank","blank.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"blank",nil)
  return
}
