package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
)

func (hnd *Handler) Edr(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"edr.html",nil)
  return
}
