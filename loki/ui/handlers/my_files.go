package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)


func (hnd *Handler) ListFiles(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("listfiles","listfiles.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"listfiles",nil)
  return
}

func (hnd *Handler) ViewFile(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("blank","blank.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"blank",nil)
  return
}
