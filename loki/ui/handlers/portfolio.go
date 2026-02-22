package handlers

import(
  "fmt"
  "net/http"
  "github.com/alphamystic/odin/lib/utils"
)


func (hnd *Handler) Createblog(res http.ResponseWriter, req *http.Request){
  _, authenticated := hnd.AuthenticateUser(res, req)
  if !authenticated {
    return // User is redirected in the helper
  }
  tpl,err := hnd.Pages.GetATemplate("creatblog","create_blog.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"creatblog",nil)
  return
}


func (hnd *Handler) Writeblog(res http.ResponseWriter, req *http.Request){
  _, authenticated := hnd.AuthenticateUser(res, req)
  if !authenticated {
    return // User is redirected in the helper
  }
  tpl,err := hnd.Pages.GetATemplate("creatblog","create_blog.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"creatblog",nil)
  return
}


func (hnd *Handler) Portfolio(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.LoadPortfolio("portfolio","portfolio.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"portfolio",nil)
  return
}


func (hnd *Handler) BlogSite(res http.ResponseWriter, req *http.Request){
  tpl,err := hnd.Pages.LoadPortfolio("blog","blog-single.tmpl")
  if err != nil {
    utils.Warning(fmt.Sprintf("%s", err))
    hnd.Internalserverror(res, req)
		return
  }
  tpl.ExecuteTemplate(res,"blog",nil)
  return
}
