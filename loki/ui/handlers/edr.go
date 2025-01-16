package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) Edr(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("edr","edr.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"edr",nil)
  return
}
