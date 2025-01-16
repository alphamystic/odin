package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) Motherships(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("motherships","motherships.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"motherships",nil)
  return
}
