package handlers

import(
  //"fmt"
  "net/http"
//  "github.com/alphamystic/odin/loki/lib/utils"
//  "github.com/alphamystic/odin/loki/lib/workers"
)

func (hnd *Handler) Odinnet(res http.ResponseWriter, req *http.Request){
  if req.Method != "POST"{
    hnd.Tpl.ExecuteTemplate(res,"blank.html",nil)
    return
  }
  hnd.Tpl.ExecuteTemplate(res,"blank.html",nil)
  return
}
