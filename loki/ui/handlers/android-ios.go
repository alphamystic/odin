package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
)

func (hnd *Handler) AndroidIOs(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"blank.html",nil)
  return
}
