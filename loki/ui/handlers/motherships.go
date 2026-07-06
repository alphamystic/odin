package handlers

import (
	"fmt"
	"net/http"
	//dfn "github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/utils"
)

func (hnd *Handler) ListMotherships(res http.ResponseWriter, req *http.Request) {
    ud, authenticated := hnd.AuthenticateUser(res, req)
    if !authenticated { return }

    token, _ := hnd.GetToken(req)
    limitStr := req.URL.Query().Get("limit")
    offsetStr := req.URL.Query().Get("offset")
    activeStr := req.URL.Query().Get("active")
    limit := utils.StringToInt(limitStr)
    offset := utils.StringToInt(offsetStr)
    if limit <= 0 {
        limit = 20
    }
    active := true
	if activeStr == "false" {
		active = false
	}
    mss, err := hnd.SRVCS.MothershipSrvs.ListMS(req.Context(), active, token, limit, offset)
    if err != nil {
        utils.Warning(fmt.Sprintf("%v",err))
        hnd.Internalserverror(res, req)
        return
    }

    tpl, err := hnd.Pages.GetATemplate("listmothership", "motherships.tmpl")
    if err != nil {
        utils.Warning(fmt.Sprintf("%v", err))
        hnd.Internalserverror(res, req)
        return
    }
    tpl.ExecuteTemplate(res, "listmothership", map[string]interface{}{
        "userdata":    ud,
        "motherships": mss,
    })
    return
}
//
// func (hnd *Handler) Motherships(res http.ResponseWriter, req *http.Request){
//   _, authenticated := hnd.AuthenticateUser(res, req)
//   if !authenticated {
//     return // User is redirected in the helper
//   }
//   tpl,err := hnd.Pages.GetATemplate("motherships","motherships.tmpl")
//   if err != nil {
//     utils.Warning(fmt.Sprintf("%s", err))
//     hnd.Internalserverror(res, req)
// 		return
//   }
//   tpl.ExecuteTemplate(res,"motherships",nil)
//   return
// }
