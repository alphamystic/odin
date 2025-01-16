package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) Blackops(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("blank","blank.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"blank",nil)
  return
}
