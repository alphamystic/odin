package definers

import (
  "github.com/dgrijalva/jwt-go"
  "github.com/alphamystic/odin/lib/utils"
)


type Mothership struct {
  OwnerID string  `json:"ownerid"`
  Name string `json:"name"`
  Password string
  MSId string `json:"mothershipid"`
  Address string `json: "address"`
  ImplantTunnel string `json:"implant_tunnel"`
  AdminTunnel string  `json:"admin_tunnel"`
  Motherships string  `json:"motherships"` // keep it as a bunch of strings and initialize it when needed
  Description string `json:"description"`
  Tls bool  `json: "tls"`
  CertPem string  `json: "cert_pem"`
  KeyPem string `json: key_pem`
  Active bool `json:"active"`
  GenCommand string `json:"generate_command"`
  Machinedata string `json:"machine_data"`
  utils.TimeStamps
}

type MaichineData struct {
  UserName string
  OsType string
  HomeDir string
  Password string
}
/*
* I am just writing down m thinking pattern here to follow along and make a decision
  So I have S/Grpc and HTTP/S
*/


type Server struct {
  Address string
  Port int //if port is 0 then the address is a url
  Protocol int
  // remeber to add the server keys just the incase of tls YOU CAN IGNORE THIS FOR NOW AS SECURE COMMS IS NOT IMPLEMENTED
  RootPem []byte // for tunneling, asigned to the client in Question
}
// @TODO there is a  descripancy on a tunnel for http/s Use admin
func GetServer(data string) (*Server,error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(utils.TOKENKEY), nil
  })
  if err != nil {
    return nil,err
  }
  var server Server
  // Extract the data from the JWT
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    data := claims["data"].([]interface{})
    server = data.(Server)
  }
  return &server,nil
}


func PrintServer(s *Server) string {
  switch s.Protocol {
  case HTTP:
    if s.Port == 0 {
      return s.Address
    } else {
       return fmt.Sprintf("http://%s:%d", s.Address, s.Port)
    }
  case HTTPS:
    if s.Port == 0 {
      return s.Address
    } else {
      return fmt.Sprintf("https://%s:%d", s.Address, s.Port)
    }
  case GRPC,SGRPC:
    return fmt.Sprintf("%s:%d", s.Address, s.Port)
  default:
    return NI
  }
  return NI
}

type ProtocolType int

const (
  HTTP ProtocolType = iota
  HTTPS
  GRPC
  SGRPC //secure grpc(tls encrypted)
  UDP
  SOCKS5
)
