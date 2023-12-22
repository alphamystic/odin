package handlers

import(
  //"fmt"
  "net/http"
//  "loki/lib/utils"
//  "loki/lib/workers"
)

func (hnd *Handler) Backdoors(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"backdoors.html",nil)
  return
}

func (hnd *Handler) Backdoorgenerator(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"blank.html",nil)
  return
}
