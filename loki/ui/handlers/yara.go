package handlers


import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) ListYaraRule(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("listfiles","listfiles.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"listfiles",nil)
  return
}

func (hnd *Handler) CreateYaraRule(res http.ResponseWriter, req *http.Request){
  /*if !IsAuthenticated(req){
    http.Redirect(res,req,"/mkubwa",http.StatusFound)//302
    return
  }*/
  if req.Method == "GET"{
    tpl,err := hnd.Pages.GetATemplate("listfiles","listfiles.tmpl")
    if err != nil {
      utils.Warning(fmt.Sprintf("%s", err))
      hnd.Internalserverror(res, req)
  		return
    }
    tpl.ExecuteTemplate(res,"listfiles",nil)
    return
  }
  if req.Method == "POST" {
    tpl,err := hnd.Pages.GetATemplate("listfiles","listfiles.tmpl")
    if err != nil {
      utils.Warning(fmt.Sprintf("%s", err))
      hnd.Internalserverror(res, req)
  		return
    }
    tpl.ExecuteTemplate(res,"listfiles",nil)
    return
  }

  http.Redirect(res,req,"/logout",http.StatusFound)
  return
}
