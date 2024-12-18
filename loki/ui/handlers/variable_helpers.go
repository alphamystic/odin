 package handlers

import (
  "fmt"
  //"log"
  "time"
  "net/http"
  "database/sql"
  "html/template"
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
  "github.com/dgrijalva/jwt-go"
//  _ "github.com/go-sql-driver/mysql"
)

type LOKI map[string]interface{}

type ErrorRes map[string]interface{}

var Registration bool

var (
  test = false
  UniversalKey = "loiuixghjpou98y7t6txcvbiuoiugyftcvbno98igtfxcfgvbioiuyft"//use this to encrypt strings/ids
)

type ErrorPage struct {
  ErrorCode int
  Data string
  Message string
  Back string
  Direction string
}

// Exposes all handlers to a db connection and the required template.
type Handler struct {
  Tpl *template.Template
  Store *http.Cookie
  Dbs *sql.DB
  RL *utils.RequestLogger
  CanWriteLogs bool
  ShutdownChan,DoneChan chan bool // channels to write into
  SRVCS *Services
}

// Initiates new handler
func NewHandler(db_connection *sql.DB, shutdownCh chan bool, doneCh chan bool,rl *utils.RequestLogger) (*Handler,error) {
  tpl,err := template.ParseGlob("./loki/ui/templates/*.html")
  if err != nil{
    utils.Warning("[-]  Failed to load templates.")
    return nil,fmt.Errorf("[-]  This is not good like: ",err)
  }
  fmt.Println("[+]  Loaded all templates.")
  utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf(" Starting LOKI server at: %s",GetCurrentTime()))
  // create db configurations
  return &Handler {
    Tpl: tpl,
    //Store: sessions.NewCookieStore([]byte(utils.RandNoLetter(30))),
    Dbs: db_connection,
    CanWriteLogs: true,
    ShutdownChan: shutdownCh,
    DoneChan: doneCh,
    RL: rl,
    SRVCS: InitializeServices(),
  },nil
}

type DateTime struct {
  Day int
  Month string
  Year int
}

func GetDateTime()(*DateTime){
  var now = time.Now()
  return &DateTime{
    Day:now.Day(),
    Month:now.Month().String(),
    Year:now.Year(),
  }
}

func GetCurrentTime() string {
  var now = time.Now()
  return now.Format("2006-01-02 15:04:05")
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


// Find a way to encrypt the tokenString
func (hnd *Handler) GetUDFromToken(req *http.Request) (*UserData,error) {
  cookie,_ := req.Cookie("Authorization")
  tokenString := cookie.Value
  // @TODO add functionality to check expiry for a jwt token and save it
  token,err := jwt.Parse(tokenString,func(tkn *jwt.Token)(interface{},error){
    if tkn.Method != jwt.SigningMethodHS256{
      return nil,fmt.Errorf("Unexepcted signing method: %v",tkn.Header["alg"])
    }
    return store,nil
  })
  if err != nil {
    return nil,fmt.Errorf("Signing error. %q",err)
  }
  if claims,ok := token.Claims.(jwt.MapClaims); ok &&  token.Valid {
    if runtimeMap,ok := claims["ud"].(map[string]interface{}); ok {
      return &UserData{
        UserId: runtimeMap["UserId"].(string),
        Admin: runtimeMap["Admin"].(bool),
      },nil
    }
  }
  return nil,dfn.NoClaimsError
}
