package handlers

import(
  //"fmt"
  "net/http"
//  "github.com/alphamystic/odin/loki/lib/utils"
//  "github.com/alphamystic/odin/loki/lib/workers"
)

func (hnd *Handler) Motherships(res http.ResponseWriter, req *http.Request){
  hnd.Tpl.ExecuteTemplate(res,"motherships.html",nil)
  return
}
