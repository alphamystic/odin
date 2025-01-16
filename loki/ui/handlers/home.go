package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)


func (hnd *Handler) Blank(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.GetATemplate("blank","blank.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"blank",nil)
  return
}



func (hnd *Handler) Internalserverror(res http.ResponseWriter, req *http.Request) {
  tpl,err := hnd.Pages.GetATemplate("error","error.tmpl")
  if err != nil{
    utils.Warning(fmt.Sprintf("%s",err))
    http.Error(res, "An error occurred", http.StatusInternalServerError)
  }
  tpl.ExecuteTemplate(res,"error",nil)
  return
}

func (hnd *Handler) Home(res http.ResponseWriter, req *http.Request){
  /*if !IsAuthenticated(req){
    http.Redirect(res,req,"/mkubwa",http.StatusFound)//302
    return
  }*/
  tpl,err := hnd.Pages.GetATemplate("home","home.tmpl")
  if err != nil{
    utils.Warning(fmt.Sprintf("%s",err))
    http.Error(res, "An error occurred", http.StatusInternalServerError)
  }
  tpl.ExecuteTemplate(res,"home",nil)
  return
}
