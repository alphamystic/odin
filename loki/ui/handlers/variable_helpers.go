package handlers

import (
  "fmt"
  "log"
  "time"
  //"database/sql"
  "html/template"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/gorilla/sessions"
//  _ "github.com/go-sql-driver/mysql"
)

type LOKI map[string]interface{}

var Registration bool

var (
  test = false
  UniversalKey = "loiuixghjpou98y7t6txcvbiuoiugyftcvbno98igtfxcfgvbioiuyft"//use this to encrypt strings/ids
)

type ErrorPage struct {
  ErrorCode int
  Message string
  Back string
}

// Exposes all handlers to a db connection and the required template.
type Handler struct {
  Tpl *template.Template
  Store *http.Cookie
  Dbs *sql.DB
  RL *utils.RequestLogger
  CanWriteLogs bool
  ShutdownChan,DoneChan chan bool // channels to write into
}

// Initiates new handler
func NewHandler(db_connection *sql.DB, shutdownCh chan bool, doneCh chan bool,rl *utils.RequestLogger) (*Handler,error) {
  tpl,err = template.ParseGlob("./templates/*.html")
  if err != nil{
    utils.Warning("[-]  Failed to load templates.")
    return nil,fmt.Errorf("[-]  This is not good like: ",err)
  }
  fmt.Println("[+]  Loaded all templates.")
  utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf(" Starting LOKI server at: %s",currentTime))
  // create db configurations
  return &Handler {
    Tpl: tpl,
    Store: sessions.NewCookieStore([]byte(utils.RandNoLetter(30))),
    Dbs: db_connection,
    CanWriteLogs: true,
    ShutdownChan: shutdownCh,
    DoneChan: doneCh,
    RL: rl,
  },nil
}

type DateTime struct {
  Day int
  Month string
  Year int
}

func GetDateTime()(*DateTime){
  return &DateTime{
    Day:now.Day(),
    Month:now.Month().String(),
    Year:now.Year(),
  }
}

type UserData struct {
  UserId string
  Admin bool
}

 var store = []byte("loki-odin")

func (hnd *Handler) GenerateJWT(ud *UserData) (string,error) {
  expTime := time.Now().Add(time.Hour * 72)
  token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
    "ud": ud,
    "exp": expTime.Unix(),
  })
  sighnedToken,err := token.SignedString(hnd.Store)
  if err != nil {
    return "",fmt.Errorf("Error signing token: %q",err)
  }
  return sighnedToken,nil
}

func (hnd *Handler) GetUDFromToken(req *http.Request) (*UserData,error) {
  session,_ := hnd.Store.Get("cookie")
  cookie,ok := session.Values["token"].(string)
  if !ok {
    return nil,dfn.UserNotLoggedIn
  }
  // @TODO add functionality to check expiry for a jwt token and save it
  token,err := jwt.Parse(cookie,func(tkn *jwt.Token)(interface{},error){
    if tkn.Method != jwt.SigningMethodHS256{
      return nil,fmt.fmt.Errorf("Unexepcted signing method: %v",tkn.Header["alg"])
    }
    return store,nil
  })
  if err != nil {
    return nil,fmt.Errorf("Signing error. %q",err)
  }
  if claims,ok := token.Claims.(jwt.MapClaims); ok &&  token.Valid {
    if runtimeMap,ok := claims["ud"].(map[string]iterface{}); ok {
      return &UserData{
        UserId: runtimeMap["UserId"],
        Admin: runtimeMap["Admin"],
      },nil
    }
  }
  return nil,dfn.NoCLaims
}
